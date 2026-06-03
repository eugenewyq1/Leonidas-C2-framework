/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_LEONIDAS_OPERATOR_TOKEN?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
