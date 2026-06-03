#!/usr/bin/env python3
"""Patch go_package in protoc-gen-go rawDesc bytes and emit valid Go string literals."""
import re
import codecs
import sys
import urllib.request
from google.protobuf import descriptor_pb2


def decode_go_concat(block: str) -> bytes:
    literals = []
    pos = block.find("=")
    for m in re.finditer(r'"(?:[^"\\]|\\.)*"', block[pos + 1 :]):
        lit = m.group(0)[1:-1]
        literals.append(codecs.decode(lit, "unicode_escape").encode("latin-1"))
    return b"".join(literals)


def extract_rawdesc_block(text: str):
    m = re.search(r"(const file_\w+_rawDesc = \"\" \+)", text)
    if not m:
        raise ValueError("rawDesc const not found")
    start = m.start()
    const_name = re.search(r"const (file_\w+_rawDesc)", text[start:]).group(1)
    # protoc-gen-go v1.36+ ends with the last string literal, then "\n\nvar ..."
    end_markers = [
        "\n\nvar (",
        "\nvar (",
        f"\n\nvar {const_name.replace('_rawDesc', '_goTypes')}",
        "\n\nvar file_",
    ]
    j = -1
    for marker in end_markers:
        j = text.find(marker, start)
        if j != -1:
            break
    if j == -1:
        j = text.find(")\n\nvar (", start)
        if j != -1:
            j += 1  # legacy format with closing paren
    if j == -1:
        raise ValueError("end of rawDesc not found")
    raw = decode_go_concat(text[start:j])
    return start, j, const_name, raw


def escape_byte(b: int) -> str:
    if b == ord("\\"):
        return "\\\\"
    if b == ord('"'):
        return '\\"'
    if b == ord("\n"):
        return "\\n"
    if b == ord("\t"):
        return "\\t"
    if b == ord("\r"):
        return "\\r"
    if 32 <= b < 127:
        return chr(b)
    return "\\x%02x" % b


def split_escaped(esc: str, width: int = 76) -> list[str]:
    chunks = []
    i = 0
    while i < len(esc):
        end = min(i + width, len(esc))
        while end > i:
            tail = esc[i:end]
            if tail.endswith("\\"):
                end -= 1
                continue
            if tail.endswith("\\x"):
                end -= 3
                continue
            lx = tail.rfind("\\x")
            if lx >= 0 and len(tail) - lx < 4:
                end -= len(tail) - lx
                continue
            break
        chunks.append(esc[i:end])
        i = end
    return chunks


def bytes_to_rawdesc_const(data: bytes, const_name: str) -> str:
    esc = "".join(escape_byte(b) for b in data)
    chunks = split_escaped(esc)
    lines = [f"const {const_name} = \"\" +"]
    for ch in chunks:
        lines.append(f'\t"{ch}" +')
    lines[-1] = lines[-1].rstrip(" +")
    return "\n".join(lines)


def patch_go_package(raw: bytes, new_pkg: str) -> bytes:
    fd = descriptor_pb2.FileDescriptorProto()
    fd.ParseFromString(raw)
    fd.options.go_package = new_pkg
    out = fd.SerializeToString()
    fd2 = descriptor_pb2.FileDescriptorProto()
    fd2.ParseFromString(out)
    if fd2.options.go_package != new_pkg:
        raise ValueError("round-trip go_package mismatch")
    return out


def rebrand_text(text: str) -> str:
    text = text.replace("github.com/bishopfox/sliver", "github.com/leonidas-c2/leonidas")
    text = text.replace("SliverRPC", "LeonidasRPC")
    text = text.replace("sliverpb.", "leonidaspb.")
    text = text.replace("protobuf/sliverpb", "protobuf/leonidaspb")
    text = text.replace('import sliverpb "', 'import leonidaspb "')
    text = text.replace(
        'sliverpb "github.com/leonidas-c2/leonidas/protobuf/leonidaspb"',
        'leonidaspb "github.com/leonidas-c2/leonidas/protobuf/leonidaspb"',
    )
    text = text.replace("package sliverpb", "package leonidaspb")
    text = text.replace("file_sliverpb_sliver_proto", "file_leonidaspb_leonidas_proto")
    text = text.replace("sliverpb/sliver.proto", "leonidaspb/leonidas.proto")
    return text


def splice_clientpb(local_path: str, bishop_url: str, new_go_package: str):
    """Keep local generated types; replace corrupted rawDesc + missing init tail from bishop."""
    local = open(local_path).read()
    bishop = urllib.request.urlopen(bishop_url).read().decode()

    ls = local.find("const file_clientpb_client_proto_rawDesc")
    if ls == -1:
        raise ValueError("local rawDesc not found")

    bs, be, const_name, raw = extract_rawdesc_block(bishop)
    new_const = bytes_to_rawdesc_const(
        patch_descriptor(raw, new_go_package, rebrand=False), const_name
    )

    bt = bishop.find("\n\nvar (", be)
    if bt == -1:
        bt = bishop.find("\n\nvar file_", be)
    if bt == -1:
        raise ValueError("bishop var block not found")

    out = local[:ls] + new_const + rebrand_text(bishop[bt:])
    open(local_path, "w").write(out)
    print("spliced", local_path)


def rebrand_descriptor(fd: descriptor_pb2.FileDescriptorProto) -> None:
    """Update module paths and package names inside a FileDescriptorProto."""
    if fd.package == "sliverpb":
        fd.package = "leonidaspb"
    if fd.name.startswith("sliverpb/"):
        fd.name = "leonidaspb/leonidas.proto"
    for i, dep in enumerate(fd.dependency):
        fd.dependency[i] = dep.replace("sliverpb/", "leonidaspb/").replace(
            "sliver.proto", "leonidas.proto"
        )
    for svc in fd.service:
        if svc.name == "SliverRPC":
            svc.name = "LeonidasRPC"
        for method in svc.method:
            if method.input_type.startswith(".sliverpb."):
                method.input_type = method.input_type.replace(".sliverpb.", ".leonidaspb.", 1)
            if method.output_type.startswith(".sliverpb."):
                method.output_type = method.output_type.replace(".sliverpb.", ".leonidaspb.", 1)


def patch_descriptor(raw: bytes, new_go_package: str, rebrand: bool) -> bytes:
    fd = descriptor_pb2.FileDescriptorProto()
    fd.ParseFromString(raw)
    fd.options.go_package = new_go_package
    if rebrand:
        rebrand_descriptor(fd)
    out = fd.SerializeToString()
    fd2 = descriptor_pb2.FileDescriptorProto()
    fd2.ParseFromString(out)
    if fd2.options.go_package != new_go_package:
        raise ValueError("go_package mismatch after patch")
    return out


def restore_from_bishop(local_path: str, bishop_url: str, new_go_package: str, rebrand: bool = False):
    bishop = urllib.request.urlopen(bishop_url).read().decode()
    start, end, const_name, raw = extract_rawdesc_block(bishop)
    new_raw = patch_descriptor(raw, new_go_package, rebrand=rebrand)
    out_const_name = rebrand_text(const_name) if rebrand else const_name
    new_const = bytes_to_rawdesc_const(new_raw, out_const_name)
    out = rebrand_text(bishop[:start]) + new_const + "\n\n" + rebrand_text(bishop[end:].lstrip("\n"))
    open(local_path, "w").write(out)
    print("restored", local_path)


if __name__ == "__main__":
    restore_from_bishop(
        "protobuf/dnspb/dns.pb.go",
        "https://raw.githubusercontent.com/BishopFox/sliver/master/protobuf/dnspb/dns.pb.go",
        "github.com/leonidas-c2/leonidas/protobuf/dnspb",
    )
    restore_from_bishop(
        "protobuf/leonidaspb/leonidas.pb.go",
        "https://raw.githubusercontent.com/BishopFox/sliver/master/protobuf/sliverpb/sliver.pb.go",
        "github.com/leonidas-c2/leonidas/protobuf/leonidaspb",
        rebrand=True,
    )
    restore_from_bishop(
        "protobuf/rpcpb/services.pb.go",
        "https://raw.githubusercontent.com/BishopFox/sliver/master/protobuf/rpcpb/services.pb.go",
        "github.com/leonidas-c2/leonidas/protobuf/rpcpb",
        rebrand=True,
    )
    splice_clientpb(
        "protobuf/clientpb/client.pb.go",
        "https://raw.githubusercontent.com/BishopFox/sliver/master/protobuf/clientpb/client.pb.go",
        "github.com/leonidas-c2/leonidas/protobuf/clientpb",
    )
