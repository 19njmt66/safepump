import { useState, useEffect, useRef } from 'react'
import { 
  useAccount, 
  useConnect, 
  useDisconnect, 
  useReadContract, 
  useWriteContract, 
  useBalance,
  usePublicClient
} from 'wagmi'
import { parseEther, formatEther } from 'viem'
import PriceChart from './components/PriceChart'
import { 
  Rocket, 
  Wallet, 
  RefreshCw, 
  Search, 
  TrendingUp, 
  ShieldCheck, 
  Coins, 
  User, 
  ArrowUpDown, 
  Check, 
  AlertCircle,
  ArrowLeft,
  Plus
} from 'lucide-react'
import { SafePumpFactoryABI, ERC20ABI } from './abi'

// Configurable constants
const DEFAULT_FACTORY_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3" // Local Anvil fallback
const API_URL = "http://localhost:8080/api/v1"
const WS_URL = "ws://localhost:8080/ws"

interface Token {
  address: string
  creator: string
  name: string
  symbol: string
  tokens_sold: string
  eth_raised: string
  migrated: boolean
  pair_address: string
  description?: string
  image_url?: string
  created_at: string
}

interface Trade {
  id: number
  token_address: string
  tx_hash: string
  block_number: number
  timestamp: string
  is_buy: boolean
  buyer_or_seller: string
  token_amount: string
  eth_amount: string
  fee: string
}

export default function App() {
  const { address, chainId, isConnected } = useAccount()
  const { connect, connectors } = useConnect()
  const { disconnect } = useDisconnect()
  const { writeContractAsync } = useWriteContract()
  const publicClient = usePublicClient()

  // App State
  const [factoryAddress, setFactoryAddress] = useState(DEFAULT_FACTORY_ADDRESS)
  const [tokens, setTokens] = useState<Token[]>([])
  const [selectedToken, setSelectedToken] = useState<Token | null>(null)
  const [trades, setTrades] = useState<Trade[]>([])
  const [searchQuery, setSearchQuery] = useState("")
  const [sortOrder, setSortOrder] = useState<"progress" | "new">("progress")
  const [activeFilterTab, setActiveFilterTab] = useState<"all" | "created">("all")
  const [lastTrade, setLastTrade] = useState<Trade | null>(null)

  // Create Token Form State
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [newTokenName, setNewTokenName] = useState("")
  const [newTokenSymbol, setNewTokenSymbol] = useState("")
  const [newTokenDescription, setNewTokenDescription] = useState("")
  const [newTokenImageUrl, setNewTokenImageUrl] = useState("")
  const [isCreatingToken, setIsCreatingToken] = useState(false)

  // Swap State
  const [isBuyTab, setIsBuyTab] = useState(true)
  const [swapAmount, setSwapAmount] = useState("")
  const [estimatedOutput, setEstimatedOutput] = useState("0")
  const [isSwapping, setIsSwapping] = useState(false)

  // Metadata edit state
  const [isEditingMetadata, setIsEditingMetadata] = useState(false)
  const [metadataDescription, setMetadataDescription] = useState("")
  const [metadataImageUrl, setMetadataImageUrl] = useState("")
  const [isUpdatingMetadata, setIsUpdatingMetadata] = useState(false)

  const socketRef = useRef<WebSocket | null>(null)

  // Get user balance
  const { data: ethBalance } = useBalance({ address, chainId })
  const formattedEthBalance = ethBalance ? formatEther(ethBalance.value) : undefined

  // Read selected token's user balance
  const { data: userTokenBalance, refetch: refetchTokenBalance } = useReadContract({
    address: selectedToken?.address as `0x${string}`,
    abi: ERC20ABI,
    functionName: 'balanceOf',
    args: address ? [address] : undefined,
    query: { enabled: !!selectedToken && !!address }
  })

  // Read selected token's Factory allowance for selling
  const { data: userAllowance, refetch: refetchAllowance } = useReadContract({
    address: selectedToken?.address as `0x${string}`,
    abi: ERC20ABI,
    functionName: 'allowance',
    args: address ? [address, factoryAddress as `0x${string}`] : undefined,
    query: { enabled: !!selectedToken && !!address }
  })

  const isCreator = !!selectedToken && !!address && selectedToken.creator.toLowerCase() === address.toLowerCase()

  // Read creator's locked amount
  const { data: creatorLockedAmount, refetch: refetchLockedAmount } = useReadContract({
    address: selectedToken?.address as `0x${string}`,
    abi: ERC20ABI,
    functionName: 'getLockedCreatorAmount',
    query: { enabled: !!selectedToken && isCreator }
  })

  // Read migration timestamp
  const { data: migrationTime } = useReadContract({
    address: selectedToken?.address as `0x${string}`,
    abi: ERC20ABI,
    functionName: 'migrationTime',
    query: { enabled: !!selectedToken && isCreator && selectedToken.migrated }
  })

  // Fetch initial token list from REST API
  const fetchTokens = async () => {
    try {
      const res = await fetch(`${API_URL}/tokens?sort=${sortOrder}`)
      const data = await res.json()
      if (Array.isArray(data)) {
        setTokens(data)
        // Set default selected token if none selected
        if (data.length > 0 && !selectedToken) {
          setSelectedToken(data[0])
        }
      }
    } catch (e) {
      console.error("Failed to fetch tokens:", e)
    }
  }

  // Fetch trades for selected token
  const fetchTrades = async (tokenAddr: string) => {
    try {
      const res = await fetch(`${API_URL}/tokens/${tokenAddr}/trades`)
      const data = await res.json()
      if (Array.isArray(data)) {
        setTrades(data)
      }
    } catch (e) {
      console.error("Failed to fetch trades:", e)
    }
  }

  useEffect(() => {
    fetchTokens()
  }, [sortOrder])

  useEffect(() => {
    if (selectedToken) {
      fetchTrades(selectedToken.address)
      refetchTokenBalance()
      refetchAllowance()
      refetchLockedAmount()
      setLastTrade(null)
      setMetadataDescription(selectedToken.description || "")
      setMetadataImageUrl(selectedToken.image_url || "")
      setIsEditingMetadata(false)
    }
  }, [selectedToken, address])

  const selectedTokenRef = useRef<Token | null>(null)
  useEffect(() => {
    selectedTokenRef.current = selectedToken
  }, [selectedToken])

  // Establish WebSockets Live Sync
  useEffect(() => {
    const connectWS = () => {
      const socket = new WebSocket(WS_URL)
      socketRef.current = socket

      socket.onopen = () => {
        console.log("Connected to Live WebSocket Stream")
      }

      socket.onmessage = (event) => {
        try {
          const packet = JSON.parse(event.data)
          if (!packet || !packet.type) return

          const currentSelected = selectedTokenRef.current

          if (packet.type === "token_created") {
            const newToken: Token = packet.data
            setTokens(prev => [newToken, ...prev])
          } else if (packet.type === "trade_processed") {
            const newTrade: Trade = packet.data
            
            // If the trade belongs to currently selected token, prepend to feed
             if (currentSelected && newTrade.token_address.toLowerCase() === currentSelected.address.toLowerCase()) {
              setTrades(prev => [newTrade, ...prev])
              setLastTrade(newTrade)
              refetchTokenBalance()
              refetchAllowance()
              refetchLockedAmount()
            }

            // Update stats on token in list
            setTokens(prev => prev.map(t => {
              if (t.address.toLowerCase() === newTrade.token_address.toLowerCase()) {
                const soldBig = BigInt(t.tokens_sold)
                const raisedBig = BigInt(t.eth_raised)
                const tradeSold = BigInt(newTrade.token_amount)
                const tradeRaised = BigInt(newTrade.eth_amount)

                if (newTrade.is_buy) {
                  return {
                    ...t,
                    tokens_sold: (soldBig + tradeSold).toString(),
                    eth_raised: (raisedBig + tradeRaised).toString()
                  }
                } else {
                  return {
                    ...t,
                    tokens_sold: (soldBig - tradeSold).toString(),
                    eth_raised: (raisedBig - tradeRaised).toString()
                  }
                }
              }
              return t
            }))

            // Update stats on selected token directly
            if (currentSelected && currentSelected.address.toLowerCase() === newTrade.token_address.toLowerCase()) {
              setSelectedToken(prev => {
                if (!prev) return null
                const soldBig = BigInt(prev.tokens_sold)
                const raisedBig = BigInt(prev.eth_raised)
                const tradeSold = BigInt(newTrade.token_amount)
                const tradeRaised = BigInt(newTrade.eth_amount)
                return {
                  ...prev,
                  tokens_sold: newTrade.is_buy ? (soldBig + tradeSold).toString() : (soldBig - tradeSold).toString(),
                  eth_raised: newTrade.is_buy ? (raisedBig + tradeRaised).toString() : (raisedBig - tradeRaised).toString()
                }
              })
            }
          } else if (packet.type === "token_migrated") {
            const migratedToken: Token = packet.data
            setTokens(prev => prev.map(t => {
              if (t.address.toLowerCase() === migratedToken.address.toLowerCase()) {
                return { ...t, migrated: true, pair_address: migratedToken.pair_address }
              }
              return t
            }))
            if (currentSelected && currentSelected.address.toLowerCase() === migratedToken.address.toLowerCase()) {
              setSelectedToken(prev => prev ? { ...prev, migrated: true, pair_address: migratedToken.pair_address } : null)
            }
          }
        } catch (e) {
          console.error("WS event parse error:", e)
        }
      }

      socket.onclose = () => {
        console.log("WS connection lost. Retrying in 3 seconds...")
        setTimeout(connectWS, 3000)
      }
    };

    connectWS()

    return () => {
      if (socketRef.current) {
        socketRef.current.close()
      }
    }
  }, [])

  // Vesting countdown timer logic for the creator
  const [vestingCountdown, setVestingCountdown] = useState<string>("")
  useEffect(() => {
    if (!selectedToken || !selectedToken.migrated || !migrationTime) {
      setVestingCountdown("")
      return
    }

    const interval = setInterval(() => {
      const now = Math.floor(Date.now() / 1000)
      const vestingStart = Number(migrationTime) + 48 * 3600

      if (now < vestingStart) {
        const diff = vestingStart - now
        const hours = Math.floor(diff / 3600)
        const minutes = Math.floor((diff % 3600) / 60)
        const seconds = diff % 60
        setVestingCountdown(`${hours}h ${minutes}m ${seconds}s`)
      } else {
        setVestingCountdown("DÉBLOQUÉ")
        clearInterval(interval)
      }
    }, 1000)

    return () => clearInterval(interval)
  }, [selectedToken, migrationTime])

  // Swap output calculation logic
  const handleAmountChange = async (val: string) => {
    setSwapAmount(val)
    if (!val || isNaN(Number(val)) || Number(val) <= 0 || !selectedToken) {
      setEstimatedOutput("0")
      return
    }

    try {
      // Dynamic estimation using local reserve emulation to speed up UX
      const sold = BigInt(selectedToken.tokens_sold)
      const raised = BigInt(selectedToken.eth_raised)
      const Tv = BigInt("941176470000000000000000000") // 941.17M tokens
      const Ev = BigInt("3000000000000000000") // 3 ETH
      const currentTokenReserves = Tv - sold
      const currentEthReserves = Ev + raised

      if (isBuyTab) {
        const ethIn = parseEther(val)
        // 1% fee subtraction
        const netEth = (ethIn * 99n) / 100n
        // Formula: (currentTokenReserves * netEth) / (currentEthReserves + netEth)
        const numerator = currentTokenReserves * netEth
        const denominator = currentEthReserves + netEth
        const out = numerator / denominator
        setEstimatedOutput(formatEther(out))
      } else {
        const tokensIn = parseEther(val)
        // Formula: (currentEthReserves * tokensIn) / (currentTokenReserves + tokensIn)
        const numerator = currentEthReserves * tokensIn
        const denominator = currentTokenReserves + tokensIn
        const ethOutNet = numerator / denominator
        // 1% fee subtraction
        const out = (ethOutNet * 99n) / 100n
        setEstimatedOutput(formatEther(out))
      }
    } catch (e) {
      console.error("Output calculation failed:", e)
    }
  }

    // Action: Launch a Token
    const handleCreateToken = async () => {
      if (!newTokenName || !newTokenSymbol) return
      setIsCreatingToken(true)

      try {
        const tx = await writeContractAsync({
          address: factoryAddress as `0x${string}`,
          abi: SafePumpFactoryABI,
          functionName: 'createToken',
          args: [newTokenName, newTokenSymbol]
        })

        console.log("Token deployment tx broadcasted:", tx)
        
        // Wait for the transaction receipt to decode the address from event logs
        if (publicClient) {
          const receipt = await publicClient.waitForTransactionReceipt({ hash: tx })
          console.log("Transaction receipt:", receipt)
          
          // Find the TokenCreated event log
          // Topic for TokenCreated(address,address,string,string): 0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4
          const log = receipt.logs.find(l => l.topics[0] === '0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4')
          if (log) {
            const tokenAddressTopic = log.topics[1]
            if (tokenAddressTopic) {
              // Convert 32-byte hex topic to 20-byte address
              const tokenAddress = '0x' + tokenAddressTopic.slice(26)
              console.log("Decoded token address:", tokenAddress)
              
              // Post metadata to the backend
              try {
                await fetch(`${API_URL}/tokens/${tokenAddress}/metadata`, {
                  method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify({
                    creator_address: address,
                    description: newTokenDescription,
                    image_url: newTokenImageUrl
                  })
                })
              } catch (metadataErr) {
                console.error("Failed to post token metadata:", metadataErr)
              }
            }
          }
        }

        setShowCreateModal(false)
        setNewTokenName("")
        setNewTokenSymbol("")
        setNewTokenDescription("")
        setNewTokenImageUrl("")
        fetchTokens()
      } catch (e) {
        console.error("Deployment failed:", e)
        alert("Erreur lors de la création du token.")
      } finally {
        setIsCreatingToken(false)
      }
    }

  // Action: Approve Token (for Selling)
  const handleApprove = async () => {
    console.log("handleApprove called", { selectedToken, swapAmount })
    if (!selectedToken || !swapAmount) {
      console.warn("handleApprove returned early: missing token or amount")
      return
    }
    setIsSwapping(true)

    try {
      const amountToApprove = parseEther(swapAmount)
      const tx = await writeContractAsync({
        address: selectedToken.address as `0x${string}`,
        abi: ERC20ABI,
        functionName: 'approve',
        args: [factoryAddress as `0x${string}`, amountToApprove]
      })
      console.log("Approve tx broadcasted:", tx)
      if (publicClient) {
        await publicClient.waitForTransactionReceipt({ hash: tx })
      }
      refetchAllowance()
    } catch (e) {
      console.error("Approval failed:", e)
      alert("Erreur lors de l'approbation.")
    } finally {
      setIsSwapping(false)
    }
  }

  // Action: Execute Swap (Buy or Sell)
  const handleSwap = async () => {
    console.log("handleSwap called", { selectedToken, swapAmount, isBuyTab })
    if (!selectedToken || !swapAmount) {
      console.warn("handleSwap returned early: missing token or amount")
      return
    }
    setIsSwapping(true)

    try {
      if (isBuyTab) {
        // Buy Action
        const ethValue = parseEther(swapAmount)
        const tx = await writeContractAsync({
          address: factoryAddress as `0x${string}`,
          abi: SafePumpFactoryABI,
          functionName: 'buy',
          args: [selectedToken.address as `0x${string}`],
          value: ethValue
        })
        console.log("Buy tx broadcasted:", tx)
      } else {
        // Sell Action
        const tokenAmount = parseEther(swapAmount)
        const tx = await writeContractAsync({
          address: factoryAddress as `0x${string}`,
          abi: SafePumpFactoryABI,
          functionName: 'sell',
          args: [selectedToken.address as `0x${string}`, tokenAmount]
        })
        console.log("Sell tx broadcasted:", tx)
      }

      console.log("Buy/Sell tx broadcasted:", tx)
      if (publicClient) {
        await publicClient.waitForTransactionReceipt({ hash: tx })
      }

      setSwapAmount("")
      setEstimatedOutput("0")
      refetchTokenBalance()
      refetchAllowance()
    } catch (e) {
      console.error("Swap failed:", e)
      alert("Erreur lors du swap.")
    } finally {
      setIsSwapping(false)
    }
  }

  // Action: Update Token Metadata (Description / Image)
  const handleUpdateMetadata = async () => {
    if (!selectedToken) return
    setIsUpdatingMetadata(true)
    try {
      const res = await fetch(`${API_URL}/tokens/${selectedToken.address}/metadata`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          creator_address: address,
          description: metadataDescription,
          image_url: metadataImageUrl
        })
      })
      if (!res.ok) {
        throw new Error("Failed to update metadata")
      }
      const updatedToken = await res.json()
      
      // Update selectedToken state locally
      setSelectedToken(updatedToken)
      
      // Update tokens list state locally
      setTokens(prev => prev.map(t => t.address.toLowerCase() === selectedToken.address.toLowerCase() ? updatedToken : t))
      
      setIsEditingMetadata(false)
      alert("Métadonnées mises à jour avec succès !")
    } catch (e) {
      console.error(e)
      alert("Erreur lors de la mise à jour des métadonnées.")
    } finally {
      setIsUpdatingMetadata(false)
    }
  }

  // Formatting helpers
  const shortenAddress = (addr: string) => `${addr.slice(0, 6)}...${addr.slice(-4)}`
  const formatBalance = (val: string | undefined) => {
    if (!val) return "0.0000"
    const cleanVal = val.replace(/,/g, '') // remove commas if present
    const parsed = parseFloat(cleanVal)
    return isNaN(parsed) ? "0.0000" : parsed.toFixed(4)
  }
  const getProgress = (token: Token) => {
    const sold = parseFloat(formatEther(BigInt(token.tokens_sold)))
    const progress = (sold / 800000000) * 100
    return Math.min(progress, 100)
  }
  const getDeterministicGradient = (address: string) => {
    const cleanAddress = address.toLowerCase().replace('0x', '');
    const color1 = cleanAddress.slice(0, 6) || '0052ff';
    const color2 = cleanAddress.slice(6, 12) || '10b981';
    const angle = parseInt(cleanAddress.slice(12, 14), 16) % 360 || 135;
    return `linear-gradient(${angle}deg, #${color1}, #${color2})`;
  };

  // Filters
  const filteredTokens = tokens.filter(t => {
    const matchesSearch = 
      t.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
      t.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.address.toLowerCase().includes(searchQuery.toLowerCase())

    if (!matchesSearch) return false

    if (activeFilterTab === "created") {
      return !!address && t.creator.toLowerCase() === address.toLowerCase()
    }

    return true
  })

  const hasSufficientAllowance = () => {
    if (isBuyTab) return true
    if (!userAllowance || !swapAmount) return false
    try {
      return (userAllowance as bigint) >= parseEther(swapAmount)
    } catch {
      return false
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      {/* HEADER NAVBAR */}
      <header className="glass-panel" style={{ margin: '20px', padding: '16px 24px', borderRadius: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <Rocket size={32} color="#0052ff" />
          <h1 style={{ fontSize: '1.6rem', color: '#fff', margin: 0 }}>SafePump</h1>
          <span style={{ fontSize: '0.75rem', background: 'rgba(0, 82, 255, 0.15)', color: '#0052ff', border: '1px solid rgba(0, 82, 255, 0.3)', padding: '3px 8px', borderRadius: '20px', fontWeight: 'bold' }}>
            Base L2
          </span>
        </div>

        {/* FACTORY ADDRESS OVERRIDE (for easy testing on Anvil/Sepolia) */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', background: 'rgba(255,255,255,0.03)', padding: '6px 12px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
          <span style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Factory:</span>
          <input 
            type="text" 
            value={factoryAddress} 
            onChange={(e) => setFactoryAddress(e.target.value)}
            style={{ background: 'transparent', border: 'none', color: '#fff', fontSize: '0.8rem', width: '130px', fontFamily: 'monospace' }}
          />
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          {isConnected && (
            <div style={{ textAlign: 'right', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
              <span style={{ fontSize: '0.85rem', color: '#fff', fontWeight: 600 }}>{shortenAddress(address as string)}</span>
              <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>
                {ethBalance ? `${formatBalance(formattedEthBalance)} ETH` : '0.0 ETH'}
              </span>
            </div>
          )}

          {isConnected ? (
            <button className="btn btn-secondary" onClick={() => disconnect()} style={{ padding: '8px 16px' }}>
              Déconnexion
            </button>
          ) : (
            <div style={{ display: 'flex', gap: '8px' }}>
              {connectors.map((connector) => (
                <button 
                  key={connector.id}
                  className="btn btn-primary" 
                  onClick={() => connect({ connector })} 
                  style={{ padding: '8px 12px', fontSize: '0.8rem' }}
                >
                  <Wallet size={14} />
                  {connector.name === 'Coinbase Wallet' ? 'Smart Wallet' : connector.name}
                </button>
              ))}
            </div>
          )}
        </div>
      </header>
      {/* CORE CONTAINER */}
      <main className="container" style={{ maxWidth: '1600px', display: 'flex', flexDirection: 'column', gap: '24px', flex: 1, paddingBottom: '40px' }}>
        {!selectedToken ? (
          <>
            {/* HERO DASHBOARD */}
            <div className="glass-panel" style={{ padding: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'radial-gradient(circle at top left, rgba(0, 82, 255, 0.05), transparent 50%), rgba(255, 255, 255, 0.01)' }}>
              <div>
                <h1 style={{ fontSize: '1.8rem', color: '#fff', fontWeight: 800, margin: 0, fontFamily: 'Outfit' }}>SafePump Launchpad</h1>
                <p style={{ fontSize: '0.9rem', color: 'var(--color-text-muted)', margin: '6px 0 0 0' }}>La plateforme de lancement équitable la plus sécurisée sur Base L2.</p>
              </div>
              <button className="btn btn-success" onClick={() => setShowCreateModal(true)} style={{ padding: '12px 24px', fontSize: '0.95rem', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '8px' }}>
                <Plus size={18} />
                Lancer un Coin
              </button>
            </div>

            {/* FILTER AND ACTION BAR */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '16px', flexWrap: 'wrap' }}>
              {/* Tabs on the Left */}
              <div style={{ display: 'flex', gap: '8px', background: 'rgba(255,255,255,0.02)', padding: '4px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                <button 
                  className={`btn ${activeFilterTab === 'all' ? 'btn-primary' : 'btn-secondary'}`}
                  onClick={() => setActiveFilterTab('all')}
                  style={{ padding: '8px 16px', fontSize: '0.85rem' }}
                >
                  Tous les Coins
                </button>
                <button 
                  className={`btn ${activeFilterTab === 'created' ? 'btn-primary' : 'btn-secondary'}`}
                  onClick={() => setActiveFilterTab('created')}
                  style={{ padding: '8px 16px', fontSize: '0.85rem' }}
                  disabled={!isConnected}
                >
                  Mes Créations
                </button>
              </div>

              {/* Search & Sort on the Right */}
              <div style={{ display: 'flex', gap: '12px', alignItems: 'center', flex: 1, justifyContent: 'flex-end', maxWidth: '600px' }}>
                <div style={{ position: 'relative', flex: 1 }}>
                  <Search size={16} color="var(--color-text-muted)" style={{ position: 'absolute', left: '12px', top: '14px' }} />
                  <input 
                    type="text" 
                    placeholder="Rechercher un token..." 
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="input-field"
                    style={{ paddingLeft: '36px' }}
                  />
                </div>

                <div style={{ display: 'flex', gap: '6px', background: 'rgba(255,255,255,0.02)', padding: '4px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                  <button 
                    className={`btn ${sortOrder === 'progress' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setSortOrder('progress')}
                    style={{ padding: '8px 12px', fontSize: '0.8rem', display: 'flex', alignItems: 'center', gap: '6px' }}
                  >
                    <TrendingUp size={14} /> Progression
                  </button>
                  <button 
                    className={`btn ${sortOrder === 'new' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setSortOrder('new')}
                    style={{ padding: '8px 12px', fontSize: '0.8rem' }}
                  >
                    Récents
                  </button>
                </div>
              </div>
            </div>

            {/* GRID OF TOKENS */}
            <div className="tokens-grid" style={{ display: 'grid', gap: '16px' }}>
              {filteredTokens.length === 0 ? (
                <div className="glass-panel" style={{ gridColumn: '1 / -1', padding: '60px 20px', textAlign: 'center', color: 'var(--color-text-muted)' }}>
                  Aucun token trouvé.
                </div>
              ) : (
                filteredTokens.map((t) => {
                  const isMyCreation = isConnected && address && t.creator.toLowerCase() === address.toLowerCase()
                  const progress = getProgress(t)
                  const ethRaisedNum = parseFloat(formatEther(BigInt(t.eth_raised)))
                  const tokensSoldNum = parseFloat(formatEther(BigInt(t.tokens_sold)))
                  const price = (3 + ethRaisedNum) / (941176470 - tokensSoldNum)
                  const mcapEth = price * 1000000000

                  return (
                    <div 
                      key={t.address}
                      onClick={() => setSelectedToken(t)}
                      className="glass-panel"
                      style={{
                        padding: '16px',
                        cursor: 'pointer',
                        transition: 'all 0.25s cubic-bezier(0.4, 0, 0.2, 1)',
                        display: 'flex',
                        flexDirection: 'column',
                        gap: '14px',
                        position: 'relative',
                        border: '1px solid var(--border-glass)',
                        overflow: 'hidden'
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.transform = 'translateY(-4px)'
                        e.currentTarget.style.boxShadow = '0 12px 30px rgba(0, 82, 255, 0.15)'
                        e.currentTarget.style.borderColor = 'rgba(0, 82, 255, 0.4)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.transform = 'translateY(0)'
                        e.currentTarget.style.boxShadow = 'none'
                        e.currentTarget.style.borderColor = 'var(--border-glass)'
                      }}
                    >
                      {isMyCreation && (
                        <div style={{ 
                          position: 'absolute', 
                          top: '8px', 
                          right: '8px', 
                          background: 'var(--color-primary-glow)', 
                          color: 'var(--color-primary)', 
                          border: '1px solid var(--color-primary)', 
                          borderRadius: '8px', 
                          padding: '1px 6px', 
                          fontSize: '0.6rem', 
                          fontWeight: 'bold',
                          zIndex: 10
                        }}>
                          👑 Créateur
                        </div>
                      )}

                      {/* Card Header with Logo Placeholder */}
                      <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
                        {t.image_url ? (
                          <img 
                            src={t.image_url} 
                            alt={t.name} 
                            style={{ width: '44px', height: '44px', borderRadius: '10px', objectFit: 'cover', flexShrink: 0, boxShadow: '0 4px 10px rgba(0, 0, 0, 0.2)' }} 
                          />
                        ) : (
                          <div style={{
                            width: '44px',
                            height: '44px',
                            borderRadius: '10px',
                            background: getDeterministicGradient(t.address),
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            fontWeight: 'bold',
                            color: '#fff',
                            fontSize: '0.9rem',
                            flexShrink: 0,
                            boxShadow: '0 4px 10px rgba(0, 0, 0, 0.2)'
                          }}>
                            {t.symbol.slice(0, 2).toUpperCase()}
                          </div>
                        )}
                        
                        <div style={{ display: 'flex', flexDirection: 'column', minWidth: 0, flex: 1 }}>
                          <h3 style={{ fontSize: '0.95rem', color: '#fff', margin: 0, fontWeight: 700, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                            {t.name}
                          </h3>
                          <div style={{ display: 'flex', gap: '6px', alignItems: 'center', marginTop: '2px' }}>
                            <span style={{ fontSize: '0.7rem', color: 'var(--color-primary)', fontWeight: 600 }}>{t.symbol}</span>
                            <span style={{ fontSize: '0.65rem', color: 'var(--color-text-dark)' }}>•</span>
                            <span style={{ fontSize: '0.65rem', color: 'var(--color-text-muted)', fontFamily: 'monospace' }}>
                              {shortenAddress(t.address)}
                            </span>
                          </div>
                        </div>
                      </div>

                      <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.7rem', color: 'var(--color-text-muted)' }}>
                          <span>Courbe :</span>
                          <span style={{ fontWeight: 'bold', color: '#fff' }}>{progress.toFixed(1)}%</span>
                        </div>
                        <div style={{ width: '100%', height: '4px', background: 'rgba(255,255,255,0.05)', borderRadius: '2px', overflow: 'hidden' }}>
                          <div style={{ width: `${progress}%`, height: '100%', background: t.migrated ? 'var(--color-success)' : 'var(--color-primary)' }} />
                        </div>
                      </div>

                      {/* Value Proposition / Description (30 chars limit) */}
                      <p style={{ 
                        fontSize: '0.7rem', 
                        color: 'var(--color-text-muted)', 
                        margin: '0', 
                        fontStyle: 'italic', 
                        whiteSpace: 'nowrap', 
                        overflow: 'hidden', 
                        textOverflow: 'ellipsis' 
                      }}>
                        {t.description ? (t.description.length > 30 ? t.description.slice(0, 30) + '...' : t.description) : "Aucune proposition de valeur."}
                      </p>

                      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '10px', fontSize: '0.75rem' }}>
                        <div>
                          <span style={{ fontSize: '0.65rem', color: 'var(--color-text-muted)' }}>MCap :</span>
                          <div style={{ color: 'var(--color-success)', fontWeight: 'bold', marginTop: '2px' }}>{mcapEth.toFixed(1)} ETH</div>
                        </div>
                        <div>
                          <span style={{ fontSize: '0.65rem', color: 'var(--color-text-muted)' }}>Levé :</span>
                          <div style={{ color: '#fff', fontWeight: 'bold', marginTop: '2px' }}>{ethRaisedNum.toFixed(2)} ETH</div>
                        </div>
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </>
        ) : (
          <>
            {/* RETURN BUTTON */}
            <button 
              className="btn btn-secondary" 
              onClick={() => setSelectedToken(null)}
              style={{ display: 'flex', alignItems: 'center', gap: '8px', width: 'fit-content', padding: '8px 16px', fontSize: '0.85rem' }}
            >
              <ArrowLeft size={16} />
              Retour aux Tokens
            </button>

            {/* SELECTED TOKEN VIEW */}
            {/* HEADER DETAILS */}
            <div className="glass-panel" style={{ padding: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', gap: '16px', alignItems: 'center' }}>
                {/* Token Logo / Avatar */}
                {selectedToken.image_url ? (
                  <img 
                    src={selectedToken.image_url} 
                    alt={selectedToken.name} 
                    style={{ width: '64px', height: '64px', borderRadius: '14px', objectFit: 'cover', flexShrink: 0, boxShadow: '0 8px 20px rgba(0, 0, 0, 0.3)' }} 
                  />
                ) : (
                  <div style={{
                    width: '64px',
                    height: '64px',
                    borderRadius: '14px',
                    background: getDeterministicGradient(selectedToken.address),
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontWeight: 'bold',
                    color: '#fff',
                    fontSize: '1.3rem',
                    flexShrink: 0,
                    boxShadow: '0 8px 20px rgba(0, 0, 0, 0.3)'
                  }}>
                    {selectedToken.symbol.slice(0, 2).toUpperCase()}
                  </div>
                )}

                <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                  <div style={{ display: 'flex', alignItems: 'baseline', gap: '10px' }}>
                    <h2 style={{ fontSize: '1.6rem', color: '#fff', margin: 0 }}>{selectedToken.name}</h2>
                    <span style={{ fontSize: '0.95rem', color: 'var(--color-primary)', fontWeight: 600 }}>[{selectedToken.symbol}]</span>
                  </div>
                  <div style={{ display: 'flex', gap: '16px', fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>
                    <span>Adresse : <span style={{ fontFamily: 'monospace' }}>{selectedToken.address}</span></span>
                    <span>Créateur : <span style={{ fontFamily: 'monospace' }}>{selectedToken.creator}</span></span>
                  </div>
                  {selectedToken.migrated && (
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--color-success)', fontSize: '0.85rem', fontWeight: 600, background: 'var(--color-success-glow)', border: '1px solid var(--color-success)', padding: '4px 10px', borderRadius: '8px', marginTop: '2px', width: 'fit-content' }}>
                      <ShieldCheck size={14} />
                      Liquidité migrée et LP brûlée à 100% sur Uniswap V2 : {shortenAddress(selectedToken.pair_address)}
                    </div>
                  )}
                </div>
              </div>
              <div style={{ textAlign: 'right' }}>
                <div style={{ fontSize: '1.8rem', color: 'var(--color-success)', fontWeight: 800 }}>
                  {parseFloat(formatEther(BigInt(selectedToken.eth_raised))).toFixed(3)} ETH
                </div>
                <div style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Levé sur 17 ETH</div>
              </div>
            </div>

            {/* METADATA & DESCRIPTION PANEL */}
            <div className="glass-panel" style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h3 style={{ fontSize: '1.1rem', color: '#fff', margin: 0, fontWeight: 700 }}>
                  À propos de {selectedToken.name}
                </h3>
                {isCreator && (
                  <button 
                    className="btn btn-secondary" 
                    onClick={() => {
                      setIsEditingMetadata(!isEditingMetadata);
                      setMetadataDescription(selectedToken.description || "");
                      setMetadataImageUrl(selectedToken.image_url || "");
                    }}
                    style={{ padding: '6px 12px', fontSize: '0.8rem' }}
                  >
                    {isEditingMetadata ? 'Annuler' : 'Modifier la Description / Image'}
                  </button>
                )}
              </div>

              {isEditingMetadata ? (
                /* EDIT METADATA FORM */
                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', background: 'rgba(255,255,255,0.01)', padding: '16px', borderRadius: '12px', border: '1px solid var(--border-glass)' }}>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                    <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)', fontWeight: 600 }}>Description du Meme Coin :</label>
                    <textarea
                      rows={3}
                      placeholder="Ajoutez une brève description de votre meme coin pour attirer les traders..."
                      value={metadataDescription}
                      onChange={(e) => setMetadataDescription(e.target.value)}
                      className="input-field"
                      style={{ resize: 'vertical', fontFamily: 'var(--font-family-body)', fontSize: '0.9rem', padding: '12px' }}
                    />
                  </div>

                  <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                    <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)', fontWeight: 600 }}>Lien de l'image (URL) :</label>
                    <input
                      type="text"
                      placeholder="Ex: https://image.com/my-meme.png"
                      value={metadataImageUrl}
                      onChange={(e) => setMetadataImageUrl(e.target.value)}
                      className="input-field"
                      style={{ fontSize: '0.9rem', padding: '10px 12px' }}
                    />
                  </div>

                  <button 
                    className="btn btn-primary" 
                    onClick={handleUpdateMetadata}
                    style={{ alignSelf: 'flex-start', padding: '10px 20px', fontSize: '0.85rem', marginTop: '4px' }}
                    disabled={isUpdatingMetadata}
                  >
                    {isUpdatingMetadata ? <RefreshCw className="animate-spin" size={16} /> : 'Sauvegarder les modifications'}
                  </button>
                </div>
              ) : (
                /* VIEW METADATA */
                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                  <p style={{ fontSize: '0.95rem', color: 'var(--color-text-main)', lineHeight: '1.6', margin: 0, whiteSpace: 'pre-wrap' }}>
                    {selectedToken.description || "Aucune description fournie pour ce meme coin. Le créateur peut en ajouter une à tout moment !"}
                  </p>
                  {selectedToken.image_url && (
                    <div style={{ display: 'flex', gap: '8px', alignItems: 'center', fontSize: '0.8rem', color: 'var(--color-text-muted)', marginTop: '4px' }}>
                      <span style={{ color: 'var(--color-success)' }}>✓</span> Image personnalisée chargée avec succès
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* PRICE CHART */}
            <PriceChart tokenAddress={selectedToken.address} lastTrade={lastTrade} />

            {/* CREATOR VESTING PANEL */}
            {isCreator && (
              <div className="glass-panel" style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <h3 style={{ fontSize: '1.05rem', color: '#fff', display: 'flex', alignItems: 'center', gap: '8px' }}>
                    <ShieldCheck size={20} color="var(--color-primary)" />
                    Allocation de Vesting du Créateur (5%)
                  </h3>
                  <span style={{ 
                    fontSize: '0.8rem', 
                    background: selectedToken.migrated ? 'var(--color-success-glow)' : 'rgba(255,255,255,0.03)', 
                    color: selectedToken.migrated ? 'var(--color-success)' : 'var(--color-text-muted)',
                    border: `1px solid ${selectedToken.migrated ? 'var(--color-success)' : 'var(--border-glass)'}`, 
                    padding: '4px 10px', 
                    borderRadius: '12px',
                    fontWeight: 'bold'
                  }}>
                    {selectedToken.migrated ? 'Migration Terminée' : 'Incubation en cours'}
                  </span>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '16px', marginTop: '4px' }}>
                  <div style={{ background: 'rgba(255,255,255,0.01)', padding: '12px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                    <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>Allocation Totale :</span>
                    <div style={{ fontSize: '1.2rem', fontWeight: 'bold', color: '#fff', marginTop: '4px' }}>50 000 000 {selectedToken.symbol}</div>
                  </div>

                  <div style={{ background: 'rgba(255,255,255,0.01)', padding: '12px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                    <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>Jetons Bloqués :</span>
                    <div style={{ fontSize: '1.2rem', fontWeight: 'bold', color: 'var(--color-danger)', marginTop: '4px' }}>
                      {creatorLockedAmount ? parseFloat(formatEther(creatorLockedAmount as bigint)).toLocaleString() : '50 000 000'} {selectedToken.symbol}
                    </div>
                  </div>

                  <div style={{ background: 'rgba(255,255,255,0.01)', padding: '12px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                    <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>Jetons Disponibles :</span>
                    <div style={{ fontSize: '1.2rem', fontWeight: 'bold', color: 'var(--color-success)', marginTop: '4px' }}>
                      {creatorLockedAmount ? (50000000 - parseFloat(formatEther(creatorLockedAmount as bigint))).toLocaleString() : '0'} {selectedToken.symbol}
                    </div>
                  </div>
                </div>

                <div style={{ background: 'rgba(255,255,255,0.02)', padding: '12px', borderRadius: '10px', border: '1px solid var(--border-glass)', fontSize: '0.85rem' }}>
                  {!selectedToken.migrated ? (
                    <div style={{ color: 'var(--color-text-muted)', display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <AlertCircle size={16} color="var(--color-primary)" />
                      Le déverrouillage commencera 48 heures après la migration de la liquidité (quand la courbe atteindra 100%).
                    </div>
                  ) : vestingCountdown && vestingCountdown !== "DÉBLOQUÉ" ? (
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: '#fff' }}>
                      <span style={{ color: 'var(--color-text-muted)' }}>Début du déverrouillage linéaire (1M/jour) dans :</span>
                      <strong style={{ color: 'var(--color-primary)', fontSize: '1rem', letterSpacing: '0.05em' }}>{vestingCountdown}</strong>
                    </div>
                  ) : (
                    <div style={{ color: 'var(--color-success)', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <Check size={16} />
                      Déblocage linéaire en cours. Vous pouvez transférer vos jetons disponibles directement depuis votre portefeuille.
                    </div>
                  )}
                </div>
              </div>
            )}
            {/* INTERACTIVE TRADE AND TRANS PANEL */}
            <div style={{ display: 'grid', gridTemplateColumns: '1.2fr 1fr', gap: '24px' }}>
              
              {/* TRADING FORM */}
              <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', height: 'fit-content' }}>
                <div style={{ display: 'flex', background: 'rgba(255,255,255,0.03)', borderRadius: '12px', padding: '4px' }}>
                  <button 
                    className={`btn ${isBuyTab ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => { setIsBuyTab(true); handleAmountChange(""); }}
                    style={{ flex: 1, padding: '10px' }}
                    disabled={selectedToken.migrated}
                  >
                    Swap (Acheter)
                  </button>
                  <button 
                    className={`btn ${!isBuyTab ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => { setIsBuyTab(false); handleAmountChange(""); }}
                    style={{ flex: 1, padding: '10px' }}
                  >
                    Swap (Vendre)
                  </button>
                </div>

                {isBuyTab ? (
                  /* BUY INTERFACE */
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.85rem' }}>
                        <span style={{ color: 'var(--color-text-muted)' }}>Montant à payer (ETH) :</span>
                        {isConnected && (
                          <span style={{ color: 'var(--color-text-muted)' }}>
                            Solde : <span style={{ color: '#fff', fontWeight: 600 }}>{formatBalance(formattedEthBalance)} ETH</span>
                          </span>
                        )}
                      </div>
                      <input 
                        type="number" 
                        step="0.01"
                        placeholder="0.0"
                        value={swapAmount}
                        onChange={(e) => handleAmountChange(e.target.value)}
                        className="input-field"
                        style={{ fontSize: '1.2rem', fontWeight: 'bold' }}
                      />
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                      <span style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Estimation reçue :</span>
                      <div style={{ fontSize: '1.2rem', fontWeight: 'bold', color: 'var(--color-success)', padding: '12px', background: 'rgba(255,255,255,0.02)', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                        {parseFloat(estimatedOutput).toLocaleString()} {selectedToken.symbol}
                      </div>
                    </div>
                  </div>
                ) : (
                  /* SELL INTERFACE */
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.85rem' }}>
                        <span style={{ color: 'var(--color-text-muted)' }}>Montant à vendre ({selectedToken.symbol}) :</span>
                        {isConnected && (
                          <span style={{ color: 'var(--color-text-muted)' }}>
                            Solde : <span style={{ color: '#fff', fontWeight: 600 }}>{userTokenBalance ? parseFloat(formatEther(userTokenBalance as bigint)).toLocaleString() : '0'} {selectedToken.symbol}</span>
                          </span>
                        )}
                      </div>
                      <input 
                        type="number" 
                        placeholder="0"
                        value={swapAmount}
                        onChange={(e) => handleAmountChange(e.target.value)}
                        className="input-field"
                        style={{ fontSize: '1.2rem', fontWeight: 'bold' }}
                      />
                      {isConnected && userTokenBalance && (
                        <div style={{ display: 'flex', gap: '8px', marginTop: '4px' }}>
                          {["25%", "50%", "75%", "100%"].map((pct) => (
                            <button
                              key={pct}
                              onClick={() => {
                                const balanceBig = userTokenBalance as bigint
                                const pctNum = parseFloat(pct) / 100
                                const amountBig = (balanceBig * BigInt(Math.floor(pctNum * 10000))) / 10000n
                                handleAmountChange(formatEther(amountBig))
                              }}
                              className="btn btn-secondary"
                              style={{ flex: 1, padding: '4px', fontSize: '0.75rem', borderRadius: '6px' }}
                            >
                              {pct}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                      <span style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Estimation reçue :</span>
                      <div style={{ fontSize: '1.2rem', fontWeight: 'bold', color: 'var(--color-success)', padding: '12px', background: 'rgba(255,255,255,0.02)', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                        {parseFloat(estimatedOutput).toFixed(6)} ETH
                      </div>
                    </div>
                  </div>
                )}

                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.8rem', color: 'var(--color-text-muted)', borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '12px' }}>
                  <span>Frais (0.5% + Min. Gaz) :</span>
                  <span>{isBuyTab ? 'Calculé' : 'Inclus'}</span>
                </div>

                {!isConnected ? (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', width: '100%' }}>
                    {connectors.map((connector) => (
                      <button 
                        key={connector.id}
                        className="btn btn-primary" 
                        onClick={() => connect({ connector })} 
                        style={{ width: '100%', padding: '10px', fontSize: '0.85rem' }}
                      >
                        <Wallet size={16} />
                        Connecter {connector.name === 'Coinbase Wallet' ? 'Smart Wallet' : connector.name}
                      </button>
                    ))}
                  </div>
                ) : !isBuyTab && !hasSufficientAllowance() ? (
                  <button 
                    className="btn btn-success" 
                    onClick={handleApprove} 
                    style={{ width: '100%' }}
                    disabled={isSwapping || !swapAmount}
                  >
                    {isSwapping ? <RefreshCw className="animate-spin" size={18} /> : 'Approuver les Tokens'}
                  </button>
                ) : (
                  <button 
                    className="btn btn-primary" 
                    onClick={handleSwap} 
                    style={{ width: '100%' }}
                    disabled={isSwapping || !swapAmount || parseFloat(swapAmount) <= 0}
                  >
                    {isSwapping ? <RefreshCw className="animate-spin" size={18} /> : (isBuyTab ? 'Swap (Acheter)' : 'Swap (Vendre)')}
                  </button>
                )}
              </div>

              {/* HISTORICAL TRADES */}
              <div className="glass-panel" style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '16px', height: '400px' }}>
                <h3 style={{ fontSize: '1rem', display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <ArrowUpDown size={16} color="var(--color-primary)" />
                  Flux de Transactions en Direct
                </h3>

                <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', overflowY: 'auto', flex: 1, paddingRight: '4px' }}>
                  {trades.length === 0 ? (
                    <div style={{ textAlign: 'center', color: 'var(--color-text-muted)', padding: '50px 0', fontSize: '0.85rem' }}>
                      Aucune transaction trouvée. Soyez le premier !
                    </div>
                  ) : (
                    trades.map((tr) => (
                      <div 
                        key={tr.id}
                        style={{
                          padding: '8px 12px',
                          background: 'rgba(255,255,255,0.01)',
                          borderRadius: '8px',
                          border: '1px solid var(--border-glass)',
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                          fontSize: '0.8rem'
                        }}
                      >
                        <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                          <span style={{ 
                            color: tr.is_buy ? 'var(--color-success)' : 'var(--color-danger)', 
                            fontWeight: 'bold',
                            background: tr.is_buy ? 'var(--color-success-glow)' : 'var(--color-danger-glow)',
                            padding: '2px 6px',
                            borderRadius: '6px',
                            fontSize: '0.7rem'
                          }}>
                            {tr.is_buy ? 'ACHAT' : 'VENTE'}
                          </span>
                          <span style={{ color: '#fff' }}>{parseFloat(formatEther(BigInt(tr.token_amount))).toLocaleString()}</span>
                          <span style={{ color: 'var(--color-text-muted)' }}>{selectedToken.symbol}</span>
                        </div>
                        <div style={{ textAlign: 'right', display: 'flex', flexDirection: 'column', gap: '2px' }}>
                          <span style={{ color: 'var(--color-success)', fontWeight: 600 }}>{parseFloat(formatEther(BigInt(tr.eth_amount))).toFixed(4)} ETH</span>
                          <span style={{ fontSize: '0.7rem', color: 'var(--color-text-dark)' }}>{shortenAddress(tr.buyer_or_seller)}</span>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </div>
            </div>
          </>
        )}
      </main>

      {/* LAUNCH TOKEN MODAL */}
      {showCreateModal && (
        <div style={{ position: 'fixed', top: 0, left: 0, width: '100%', height: '100%', background: 'rgba(2, 6, 23, 0.8)', backdropFilter: 'blur(8px)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
          <div className="glass-panel" style={{ width: '400px', padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
            <h3 style={{ fontSize: '1.3rem', color: '#fff', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Rocket size={22} color="var(--color-success)" />
              Lancer votre Meme Coin
            </h3>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Nom du Token :</label>
              <input 
                type="text" 
                placeholder="Ex: Safe Dog"
                value={newTokenName}
                onChange={(e) => setNewTokenName(e.target.value)}
                className="input-field"
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Symbole du Token :</label>
              <input 
                type="text" 
                placeholder="Ex: SDOG"
                value={newTokenSymbol}
                onChange={(e) => setNewTokenSymbol(e.target.value)}
                className="input-field"
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Proposition de Valeur :</label>
              <textarea 
                placeholder="Expliquez la proposition de valeur ou le mème..."
                value={newTokenDescription}
                onChange={(e) => setNewTokenDescription(e.target.value)}
                className="input-field"
                rows={2}
                style={{ resize: 'none', padding: '10px 12px' }}
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Lien de l'image (URL) :</label>
              <input 
                type="text" 
                placeholder="Ex: https://i.imgur.com/xyz.png"
                value={newTokenImageUrl}
                onChange={(e) => setNewTokenImageUrl(e.target.value)}
                className="input-field"
              />
            </div>

            <div style={{ display: 'flex', gap: '12px', marginTop: '8px' }}>
              <button 
                className="btn btn-secondary" 
                onClick={() => {
                  setShowCreateModal(false);
                  setNewTokenName("");
                  setNewTokenSymbol("");
                  setNewTokenDescription("");
                  setNewTokenImageUrl("");
                }} 
                style={{ flex: 1 }}
              >
                Annuler
              </button>
              <button 
                className="btn btn-primary" 
                onClick={handleCreateToken}
                style={{ flex: 1 }}
                disabled={isCreatingToken || !newTokenName || !newTokenSymbol}
              >
                {isCreatingToken ? <RefreshCw className="animate-spin" size={18} /> : 'Déployer'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
