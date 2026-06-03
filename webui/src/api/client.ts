// Leonidas C2 gRPC-Web API client
// Sends serialised protobuf requests to the gRPC-Web proxy at /rpcpb.LeonidasRPC/<Method>

const BASE = '/rpcpb.LeonidasRPC'

// Operator token for dev (same value leonidas-client uses as Bearer). Create webui/.env.local:
//   VITE_LEONIDAS_OPERATOR_TOKEN=<paste operator JWT>
const operatorToken = import.meta.env.VITE_LEONIDAS_OPERATOR_TOKEN ?? ''

// grpcWebCall - Makes a gRPC-Web (binary) request and returns the raw response bytes.
// In production, the request body is a serialised protobuf message; for now we use
// JSON-in-bytes as a placeholder until protoc-gen-grpc-web stubs are generated.
async function grpcWebCall<Req, Resp>(method: string, req: Req): Promise<Resp> {
  const encoder = new TextEncoder()
  const reqJson = JSON.stringify(req)
  const reqBytes = encoder.encode(reqJson)

  // Build gRPC-Web frame: 1 byte flags (0 = no compression) + 4 bytes big-endian length
  const frame = new Uint8Array(5 + reqBytes.length)
  frame[0] = 0
  new DataView(frame.buffer).setUint32(1, reqBytes.length, false)
  frame.set(reqBytes, 5)

  const response = await fetch(`${BASE}/${method}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/grpc-web+proto',
      'X-Grpc-Web': '1',
      ...(operatorToken ? { Authorization: `Bearer ${operatorToken}` } : {}),
    },
    body: frame,
  })

  if (!response.ok) {
    throw new Error(`gRPC call ${method} failed: ${response.status}`)
  }

  const buf = await response.arrayBuffer()
  // Strip the 5-byte gRPC data frame header
  const data = new Uint8Array(buf).slice(5)
  const decoder = new TextDecoder()
  return JSON.parse(decoder.decode(data)) as Resp
}

// ---------------------------------------------------------------------------
// API surface — mirrors the LeonidasRPC service
// ---------------------------------------------------------------------------

export interface Session {
  ID: string
  Name: string
  RemoteAddress: string
  Hostname: string
  OS: string
  Arch: string
  PID: number
  LastCheckin: string
  ActiveC2: string
}

export interface Beacon {
  ID: string
  Name: string
  RemoteAddress: string
  Hostname: string
  OS: string
  Interval: number
  Jitter: number
  NextCheckin: string
  ActiveC2: string
}

export interface Job {
  ID: number
  Name: string
  Description: string
  Protocol: string
  Port: number
}

export interface Operator {
  Name: string
  Online: boolean
}

export interface LootItem {
  ID: string
  Name: string
  Type: string
  CreatedAt: string
}

export const api = {
  getSessions: () => grpcWebCall<Record<string, never>, { Sessions: Session[] }>('GetSessions', {}),
  getBeacons:  () => grpcWebCall<Record<string, never>, { Beacons: Beacon[] }>('GetBeacons', {}),
  getJobs:     () => grpcWebCall<Record<string, never>, { Active: Job[] }>('GetJobs', {}),
  getOperators: () => grpcWebCall<Record<string, never>, { Operators: Operator[] }>('GetOperators', {}),
  getLootFiles: () => grpcWebCall<Record<string, never>, { LootFiles: LootItem[] }>('LootFiles', {}),
  getVersion:  () => grpcWebCall<Record<string, never>, { Major: number; Minor: number; Patch: number; Commit: string }>('GetVersion', {}),

  startMTLSListener: (host: string, port: number) =>
    grpcWebCall<{ Host: string; Port: number }, { JobID: number }>('StartMTLSListener', { Host: host, Port: port }),

  startICMPListener: (host: string) =>
    grpcWebCall<{ Host: string }, { JobID: number }>('StartICMPListener', { Host: host }),

  killJob: (id: number) =>
    grpcWebCall<{ ID: number }, Record<string, never>>('KillJob', { ID: id }),
}
