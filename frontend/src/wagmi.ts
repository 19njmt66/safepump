import { http, createConfig } from 'wagmi'
import { foundry, localhost } from 'wagmi/chains'
import { injected, coinbaseWallet } from 'wagmi/connectors'

export const config = createConfig({
  chains: [foundry, localhost],
  connectors: [
    injected(),
    coinbaseWallet({ appName: 'SafePump' }),
  ],
  transports: {
    [foundry.id]: http(`http://${typeof window !== 'undefined' ? window.location.hostname : '127.0.0.1'}:8545`),
    [localhost.id]: http(`http://${typeof window !== 'undefined' ? window.location.hostname : '127.0.0.1'}:8545`),
  },
})

declare module 'wagmi' {
  interface Register {
    config: typeof config
  }
}
