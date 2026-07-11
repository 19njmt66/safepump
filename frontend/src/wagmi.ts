import { http, createConfig } from 'wagmi'
import { foundry, localhost } from 'wagmi/chains'
import { injected } from 'wagmi/connectors'

export const config = createConfig({
  chains: [foundry, localhost],
  connectors: [
    injected(),
  ],
  transports: {
    [foundry.id]: http('http://127.0.0.1:8545'),
    [localhost.id]: http('http://127.0.0.1:8545'),
  },
})

declare module 'wagmi' {
  interface Register {
    config: typeof config
  }
}
