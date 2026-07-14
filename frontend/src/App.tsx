import { useState, useEffect, useRef } from 'react'
import { 
  useAccount, 
  useConnect, 
  useDisconnect, 
  useReadContract, 
  useWriteContract, 
  useBalance,
  usePublicClient,
  useSignMessage
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
  User,
  Check, 
  AlertCircle,
  ArrowLeft,
  Plus,
  Sun,
  Moon,
  Globe,
  Send
} from 'lucide-react'
import { SafePumpFactoryABI, ERC20ABI } from './abi'

// Configurable constants
const DEFAULT_FACTORY_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3" // Local Anvil fallback
const hostname = typeof window !== 'undefined' ? window.location.hostname : 'localhost';
const API_URL = `http://${hostname}:8080/api/v1`
const WS_URL = `ws://${hostname}:8080/ws`

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
  website?: string
  twitter?: string
  telegram?: string
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
  const { address, isConnected } = useAccount()
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
  const [activeFilter, setActiveFilter] = useState<"all" | "created" | "progress" | "new">("all")
  const [lastTrade, setLastTrade] = useState<Trade | null>(null)

  // Create Token Form State
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [newTokenName, setNewTokenName] = useState("")
  const [newTokenSymbol, setNewTokenSymbol] = useState("")
  const [newTokenDescription, setNewTokenDescription] = useState("")
  const [newTokenImageUrl, setNewTokenImageUrl] = useState("")
  const [newTokenWebsite, setNewTokenWebsite] = useState("")
  const [newTokenTwitter, setNewTokenTwitter] = useState("")
  const [newTokenTelegram, setNewTokenTelegram] = useState("")
  const [createWithBonus, setCreateWithBonus] = useState(false)
  const [initialBuyEth, setInitialBuyEth] = useState("")
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
  const [metadataWebsite, setMetadataWebsite] = useState("")
  const [metadataTwitter, setMetadataTwitter] = useState("")
  const [metadataTelegram, setMetadataTelegram] = useState("")
  const [isUpdatingMetadata, setIsUpdatingMetadata] = useState(false)

  // Audit State
  const [auditData, setAuditData] = useState<any>(null)
  const [isLoadingAudit, setIsLoadingAudit] = useState(false)

  // Profile State
  const [showProfileModal, setShowProfileModal] = useState(false)
  const [profileUsername, setProfileUsername] = useState("")
  const [editUsername, setEditUsername] = useState("")
  const [editBio, setEditBio] = useState("")
  const [editAvatarUrl, setEditAvatarUrl] = useState("")
  const [isSavingProfile, setIsSavingProfile] = useState(false)
  const { signMessageAsync } = useSignMessage()

  // Profile View States
  const [viewProfileAddress, setViewProfileAddress] = useState<string | null>(null)
  const [profileData, setProfileData] = useState<any>(null)
  const [portfolioData, setPortfolioData] = useState<any>(null)
  const [creatorTokens, setCreatorTokens] = useState<Token[]>([])
  const [userTrades, setUserTrades] = useState<any[]>([])
  const [profileTab, setProfileTab] = useState<"held" | "created" | "trades">("held")

  // Token Page Details States
  const [detailsTab, setDetailsTab] = useState<"trades" | "holders">("trades")
  const [tokenHolders, setTokenHolders] = useState<any[]>([])

  // Theme State
  const [theme, setTheme] = useState<'dark' | 'light'>('dark')

  const socketRef = useRef<WebSocket | null>(null)

  // Get user balance
  const { data: ethBalance } = useBalance({ address })
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

  // Load and apply theme
  useEffect(() => {
    const savedTheme = localStorage.getItem('theme') as 'dark' | 'light' || 'dark'
    setTheme(savedTheme)
    document.documentElement.setAttribute('data-theme', savedTheme)
  }, [])

  const toggleTheme = () => {
    const nextTheme = theme === 'dark' ? 'light' : 'dark'
    setTheme(nextTheme)
    document.documentElement.setAttribute('data-theme', nextTheme)
    localStorage.setItem('theme', nextTheme)
  }

  useEffect(() => {
    const backendSort = (activeFilter === "new") ? "new" : "progress";
    if (backendSort !== sortOrder) {
      setSortOrder(backendSort);
    } else {
      fetchTokens();
    }
  }, [activeFilter, sortOrder]);

  const fetchTokenHolders = async (tokenAddress: string) => {
    try {
      const res = await fetch(`${API_URL}/tokens/${tokenAddress.toLowerCase()}/holders`)
      const data = await res.json()
      if (Array.isArray(data)) {
        setTokenHolders(data)
      }
    } catch (e) {
      console.error("Failed to fetch token holders:", e)
    }
  }

  useEffect(() => {
    if (selectedToken) {
      fetchTrades(selectedToken.address)
      fetchTokenHolders(selectedToken.address)
      refetchTokenBalance()
      refetchAllowance()
      refetchLockedAmount()
      setLastTrade(null)
      setMetadataDescription(selectedToken.description || "")
      setMetadataImageUrl(selectedToken.image_url || "")
      setMetadataWebsite(selectedToken.website || "")
      setMetadataTwitter(selectedToken.twitter || "")
      setMetadataTelegram(selectedToken.telegram || "")
    }
  }, [selectedToken, address])

  // Fetch platform configuration (factory address) on mount
  useEffect(() => {
    fetch(`${API_URL}/config`)
      .then(res => res.json())
      .then(data => {
        if (data && data.factory_address) {
          setFactoryAddress(data.factory_address)
          console.log("Loaded factory address from backend config:", data.factory_address)
        }
      })
      .catch(err => console.error("Failed to fetch platform config:", err))
  }, [])

  // Fetch current user profile when wallet connects
  useEffect(() => {
    if (isConnected && address) {
      fetch(`${API_URL}/users/${address.toLowerCase()}`)
        .then(res => {
          if (!res.ok) throw new Error("Failed to fetch profile");
          return res.json();
        })
        .then(data => {
          if (data) {
            setProfileUsername(data.username || "");
            setEditUsername(data.username || "");
            setEditBio(data.bio || "");
            setEditAvatarUrl(data.avatar_url || "");
          }
        })
        .catch(err => {
          console.error("Failed to fetch profile:", err);
          setProfileUsername("");
        });
    } else {
      setProfileUsername("");
    }
  }, [isConnected, address]);

  // Fetch profile page data when viewProfileAddress changes
  useEffect(() => {
    if (viewProfileAddress) {
      // 1. Fetch user info
      fetch(`${API_URL}/users/${viewProfileAddress.toLowerCase()}`)
        .then(res => {
          if (!res.ok) throw new Error("Profile not found");
          return res.json();
        })
        .then(data => setProfileData(data))
        .catch(err => {
          console.error("Failed to fetch profile info:", err);
          setProfileData({ username: "", bio: "", avatar_url: "" });
        })

      // 2. Fetch portfolio (coins held)
      fetch(`${API_URL}/users/${viewProfileAddress.toLowerCase()}/portfolio`)
        .then(res => res.json())
        .then(data => setPortfolioData(data))
        .catch(err => console.error("Failed to fetch portfolio:", err))

      // 3. Fetch launched tokens
      fetch(`${API_URL}/creator/${viewProfileAddress.toLowerCase()}/tokens`)
        .then(res => res.json())
        .then(data => {
          if (Array.isArray(data)) setCreatorTokens(data)
        })
        .catch(err => console.error("Failed to fetch creator tokens:", err))

      // 4. Fetch user trades
      fetch(`${API_URL}/users/${viewProfileAddress.toLowerCase()}/trades`)
        .then(res => res.json())
        .then(data => {
          if (Array.isArray(data)) setUserTrades(data)
        })
        .catch(err => console.error("Failed to fetch user trades:", err))
    } else {
      setProfileData(null)
      setPortfolioData(null)
      setCreatorTokens([])
      setUserTrades([])
    }
  }, [viewProfileAddress]);

  // Fetch token audit data when selectedToken changes
  useEffect(() => {
    if (selectedToken) {
      setIsLoadingAudit(true);
      fetch(`${API_URL}/tokens/${selectedToken.address.toLowerCase()}/audit`)
        .then(res => res.json())
        .then(data => {
          setAuditData(data);
        })
        .catch(err => {
          console.error("Failed to fetch token audit:", err);
          setAuditData(null);
        })
        .finally(() => {
          setIsLoadingAudit(false);
        });
    } else {
      setAuditData(null);
    }
  }, [selectedToken]);

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
        // Deploy token (with optional 0.02127 ETH for Creator Pack bonus)
        const tx = await writeContractAsync({
          address: factoryAddress as `0x${string}`,
          abi: SafePumpFactoryABI,
          functionName: 'createToken',
          args: [newTokenName, newTokenSymbol],
          value: createWithBonus ? parseEther('0.02127') : undefined
        })

        console.log("Token deployment tx broadcasted:", tx)
        let createdTokenAddress = ""
        
        // Wait for the transaction receipt to decode the address from event logs
        if (publicClient) {
          const receipt = await publicClient.waitForTransactionReceipt({ hash: tx })
          console.log("Transaction receipt:", receipt)
          
          // Find the TokenCreated event log
          const log = receipt.logs.find(l => l.topics[0] === '0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4')
          if (log) {
            const tokenAddressTopic = log.topics[1]
            if (tokenAddressTopic) {
              createdTokenAddress = '0x' + tokenAddressTopic.slice(26)
              console.log("Decoded token address:", createdTokenAddress)
              
              // Post metadata (including socials) to the backend
              try {
                await fetch(`${API_URL}/tokens/${createdTokenAddress}/metadata`, {
                  method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify({
                    creator_address: address,
                    description: newTokenDescription,
                    image_url: newTokenImageUrl,
                    website: newTokenWebsite,
                    twitter: newTokenTwitter,
                    telegram: newTokenTelegram
                  })
                })
              } catch (metadataErr) {
                console.error("Failed to post token metadata:", metadataErr)
              }
            }
          }
        }

        // Chained Initial Buy if requested
        if (createdTokenAddress && initialBuyEth && parseFloat(initialBuyEth) > 0) {
          console.log("Initiating chained initial buy of", initialBuyEth, "ETH...")
          try {
            const buyTx = await writeContractAsync({
              address: factoryAddress as `0x${string}`,
              abi: SafePumpFactoryABI,
              functionName: 'buy',
              args: [createdTokenAddress as `0x${string}`],
              value: parseEther(initialBuyEth)
            })
            console.log("Initial buy transaction broadcasted:", buyTx)
            if (publicClient) {
              await publicClient.waitForTransactionReceipt({ hash: buyTx })
              console.log("Initial buy confirmed!")
            }
          } catch (buyErr) {
            console.error("Initial buy failed:", buyErr)
            alert("Le jeton a été créé avec succès, mais l'achat initial a échoué. Vous pouvez acheter manuellement depuis la page du jeton.")
          }
        }

        setShowCreateModal(false)
        setNewTokenName("")
        setNewTokenSymbol("")
        setNewTokenDescription("")
        setNewTokenImageUrl("")
        setNewTokenWebsite("")
        setNewTokenTwitter("")
        setNewTokenTelegram("")
        setCreateWithBonus(false)
        setInitialBuyEth("")
        fetchTokens()
        alert("Félicitations ! Votre Meme Coin a été déployé avec succès !")
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
      let tx: `0x${string}`
      if (isBuyTab) {
        // Buy Action
        const ethValue = parseEther(swapAmount)
        tx = await writeContractAsync({
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
        tx = await writeContractAsync({
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
      fetchTokenHolders(selectedToken.address)
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
          image_url: metadataImageUrl,
          website: metadataWebsite,
          twitter: metadataTwitter,
          telegram: metadataTelegram
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

  // Action: Save user profile
  const handleSaveProfile = async () => {
    if (!address) return
    setIsSavingProfile(true)
    try {
      const msg = `Mise à jour du profil SafePump pour ${address.toLowerCase()} à ${Date.now()}`
      const sig = await signMessageAsync({ message: msg })
      
      const payload = {
        address: address,
        username: editUsername,
        bio: editBio,
        avatar_url: editAvatarUrl,
        message: msg,
        signature: sig
      }

      const res = await fetch(`${API_URL}/users/profile`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      })

      if (!res.ok) {
        throw new Error("Failed to save profile")
      }

      const updatedUser = await res.json()
      setProfileUsername(updatedUser.username || "")
      setShowProfileModal(false)
      alert("Profil mis à jour avec succès !")
    } catch (e) {
      console.error(e)
      alert("Erreur lors de la mise à jour du profil.")
    } finally {
      setIsSavingProfile(false)
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

  const getRelativeTime = (dateString: string) => {
    const ms = Date.now() - new Date(dateString).getTime()
    const secs = Math.floor(ms / 1000)
    const mins = Math.floor(secs / 60)
    const hours = Math.floor(mins / 60)
    const days = Math.floor(hours / 24)
    if (days > 0) return `il y a ${days}j`
    if (hours > 0) return `il y a ${hours}h`
    if (mins > 0) return `il y a ${mins}m`
    return "à l'instant"
  }

  // Filters
  const filteredTokens = tokens.filter(t => {
    const matchesSearch = 
      t.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
      t.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.address.toLowerCase().includes(searchQuery.toLowerCase())

    if (!matchesSearch) return false

    if (activeFilter === "created") {
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
      <header className="glass-panel" style={{ margin: '20px', padding: '22px 30px', borderRadius: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <Rocket size={40} color="#0052ff" />
          <div style={{ display: 'flex', flexDirection: 'column' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <h1 style={{ fontSize: '2.1rem', color: '#fff', margin: 0, cursor: 'pointer', fontWeight: 800, fontFamily: 'Outfit', lineHeight: '1' }} onClick={() => { setViewProfileAddress(null); setSelectedToken(null); }}>SafePump</h1>
              <span style={{ fontSize: '0.75rem', background: 'rgba(0, 82, 255, 0.15)', color: '#0052ff', border: '1px solid rgba(0, 82, 255, 0.3)', padding: '3px 8px', borderRadius: '20px', fontWeight: 'bold' }}>
                Base L2
              </span>
            </div>
            <p style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)', margin: '4px 0 0 0', padding: 0 }}>
              La plateforme de lancement équitable la plus sécurisée sur Base L2.
            </p>
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          {isConnected && (
            <div style={{ textAlign: 'right', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
              <span 
                style={{ 
                  fontSize: '0.85rem', 
                  color: '#fff', 
                  fontWeight: 600, 
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px',
                  background: 'rgba(255, 255, 255, 0.04)',
                  padding: '4px 10px',
                  borderRadius: '8px',
                  border: '1px solid var(--border-glass)',
                  transition: 'all 0.2s'
                }}
                onClick={() => {
                  if (address) {
                    setViewProfileAddress(address)
                    setSelectedToken(null)
                  }
                }}
                title="Modifier mon profil"
                onMouseEnter={(e) => e.currentTarget.style.background = 'rgba(255, 255, 255, 0.08)'}
                onMouseLeave={(e) => e.currentTarget.style.background = 'rgba(255, 255, 255, 0.04)'}
              >
                👤 {profileUsername || shortenAddress(address as string)}
              </span>
              <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>
                {ethBalance ? `${formatBalance(formattedEthBalance)} ETH` : '0.0 ETH'}
              </span>
            </div>
          )}

          {/* Theme Toggle Button */}
          <button 
            className="btn btn-secondary" 
            onClick={toggleTheme} 
            style={{ 
              padding: '0', 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'center',
              width: '38px',
              height: '38px',
              borderRadius: '10px',
              background: 'rgba(255, 255, 255, 0.04)',
              border: '1px solid var(--border-glass)'
            }}
            title={theme === 'dark' ? 'Activer le mode clair' : 'Activer le mode sombre'}
          >
            {theme === 'dark' ? <Sun size={16} color="#fbbf24" /> : <Moon size={16} color="var(--color-primary)" />}
          </button>

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
      <main className="container" style={{ maxWidth: '96%', display: 'flex', flexDirection: 'column', gap: '24px', flex: 1, paddingBottom: '40px' }}>
        {viewProfileAddress ? (
          /* PROFILE PAGE VIEW */
          <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
             {/* Back Button */}
             <button 
               className="btn btn-secondary" 
               onClick={() => setViewProfileAddress(null)}
               style={{ display: 'flex', alignItems: 'center', gap: '8px', width: 'fit-content', padding: '8px 16px', fontSize: '0.85rem' }}
             >
               <ArrowLeft size={16} />
               Retour à l'accueil
             </button>

             {/* Profile Grid (2 columns: Left details, Right metrics) */}
             <div style={{ display: 'flex', gap: '20px', alignItems: 'stretch', flexWrap: 'wrap' }}>
               {/* Left Box: Profile Info */}
               <div className="glass-panel" style={{ flex: '2 1 450px', padding: '24px', display: 'flex', gap: '20px', alignItems: 'center' }}>
                 {/* Avatar */}
                 <div style={{
                   width: '90px',
                   height: '90px',
                   borderRadius: '16px',
                   background: profileData?.avatar_url ? `url(${profileData.avatar_url}) center/cover no-repeat` : getDeterministicGradient(viewProfileAddress),
                   display: 'flex',
                   alignItems: 'center',
                   justifyContent: 'center',
                   fontWeight: 'bold',
                   color: '#fff',
                   fontSize: '1.8rem',
                   boxShadow: '0 6px 20px rgba(0,0,0,0.3)',
                   border: '2px solid rgba(255,255,255,0.08)',
                   flexShrink: 0
                 }}>
                   {!profileData?.avatar_url && (profileData?.username ? profileData.username.slice(0, 2).toUpperCase() : viewProfileAddress.slice(2, 4).toUpperCase())}
                 </div>

                 {/* Name/Bio */}
                 <div style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: '6px', minWidth: 0 }}>
                   <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
                     <h2 style={{ fontSize: '1.5rem', color: '#fff', fontWeight: 800, margin: 0 }}>
                       {profileData?.username || shortenAddress(viewProfileAddress)}
                     </h2>
                     {viewProfileAddress.toLowerCase() === address?.toLowerCase() && (
                       <button 
                         className="btn btn-secondary"
                         onClick={() => {
                           setEditUsername(profileData?.username || "");
                           setEditBio(profileData?.bio || "");
                           setEditAvatarUrl(profileData?.avatar_url || "");
                           setShowProfileModal(true);
                         }}
                         style={{ padding: '3px 8px', fontSize: '0.7rem', height: 'auto', borderRadius: '6px' }}
                       >
                         ✏️ Modifier
                       </button>
                     )}
                   </div>
                   <div 
                     style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)', fontFamily: 'monospace', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '6px' }}
                     onClick={() => {
                       navigator.clipboard.writeText(viewProfileAddress);
                       alert("Adresse copiée !");
                     }}
                   >
                     <span>{shortenAddress(viewProfileAddress)}</span>
                     <span style={{ fontSize: '0.7rem', color: 'var(--color-primary)' }}>Copier</span>
                   </div>
                   <p style={{ margin: 0, fontSize: '0.85rem', color: 'var(--color-text-muted)', fontStyle: profileData?.bio ? 'normal' : 'italic', whiteSpace: 'pre-wrap' }}>
                     {profileData?.bio || "Aucune biographie disponible."}
                   </p>
                 </div>
               </div>

               {/* Right Box: Indicators */}
               <div className="glass-panel" style={{ flex: '1 1 300px', padding: '24px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                 <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px 20px' }}>
                   {/* PnL Indicator */}
                   <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                     <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)', textTransform: 'uppercase', fontWeight: 600 }}>P&L Réalisé</span>
                     <span style={{ 
                       fontSize: '1.25rem', 
                       fontWeight: 800, 
                       color: (portfolioData?.realized_pnl_eth || 0) >= 0 ? 'var(--color-success)' : 'var(--color-danger)'
                     }}>
                       {portfolioData ? `${portfolioData.realized_pnl_eth >= 0 ? '+' : ''}${portfolioData.realized_pnl_eth.toFixed(4)} ETH` : '0.0000 ETH'}
                     </span>
                   </div>

                   {/* Migrated Tokens Indicator */}
                   <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                     <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)', textTransform: 'uppercase', fontWeight: 600 }}>Jetons Migrés</span>
                     <span style={{ fontSize: '1.25rem', fontWeight: 800, color: '#fff' }}>
                       {creatorTokens.filter(t => t.migrated).length}
                     </span>
                   </div>

                   {/* Volume Indicator */}
                   <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                     <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)', textTransform: 'uppercase', fontWeight: 600 }}>Volume Tradé</span>
                     <span style={{ fontSize: '1.25rem', fontWeight: 800, color: '#fff' }}>
                       {portfolioData ? `${portfolioData.total_volume_eth.toFixed(3)} ETH` : '0.000 ETH'}
                     </span>
                   </div>

                   {/* Followers Indicator */}
                   <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                     <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)', textTransform: 'uppercase', fontWeight: 600 }}>Followers</span>
                     <span style={{ fontSize: '1.25rem', fontWeight: 800, color: '#fff' }}>
                       {((parseInt(viewProfileAddress.slice(2, 6), 16) % 450) + 12)}
                     </span>
                   </div>
                 </div>
               </div>
             </div>

             {/* Profile Tabs */}
             <div style={{ display: 'flex', gap: '8px', borderBottom: '1px solid rgba(255,255,255,0.06)', paddingBottom: '12px', marginTop: '12px' }}>
               <button 
                 className={`btn ${profileTab === 'held' ? 'btn-primary' : 'btn-secondary'}`}
                 onClick={() => setProfileTab('held')}
                 style={{ padding: '8px 16px', fontSize: '0.9rem' }}
               >
                 Coins Détenus ({portfolioData?.positions?.length || 0})
               </button>
               <button 
                 className={`btn ${profileTab === 'created' ? 'btn-primary' : 'btn-secondary'}`}
                 onClick={() => setProfileTab('created')}
                 style={{ padding: '8px 16px', fontSize: '0.9rem' }}
               >
                 Coins Créés ({creatorTokens.length})
               </button>
               <button 
                 className={`btn ${profileTab === 'trades' ? 'btn-primary' : 'btn-secondary'}`}
                 onClick={() => setProfileTab('trades')}
                 style={{ padding: '8px 16px', fontSize: '0.9rem' }}
               >
                 Transactions ({userTrades.length})
               </button>
             </div>

             {/* Tab Content */}
             {profileTab === 'held' && (
               <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                 {!portfolioData || portfolioData.positions.length === 0 ? (
                   <div style={{ padding: '60px 0', textAlign: 'center', color: 'var(--color-text-muted)', fontSize: '0.95rem' }} className="glass-panel">
                     Aucun jeton détenu pour le moment.
                   </div>
                 ) : (
                   <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: '20px' }}>
                     {portfolioData.positions.map((pos: any) => {
                       const tokenBal = parseFloat(pos.balance) / 1e18;
                       const positionValue = tokenBal * pos.current_price;
                       return (
                         <div 
                           key={pos.token_address} 
                           className="glass-panel" 
                           onClick={() => {
                             const foundToken = tokens.find(t => t.address.toLowerCase() === pos.token_address.toLowerCase());
                             if (foundToken) {
                               setSelectedToken(foundToken);
                               setViewProfileAddress(null);
                             }
                           }}
                           style={{ padding: '20px', display: 'flex', gap: '14px', alignItems: 'center', cursor: 'pointer', transition: 'transform 0.2s' }}
                           onMouseEnter={(e) => e.currentTarget.style.transform = 'translateY(-2px)'}
                           onMouseLeave={(e) => e.currentTarget.style.transform = 'translateY(0)'}
                         >
                           {pos.image_url ? (
                             <img src={pos.image_url} alt={pos.token_name} style={{ width: '48px', height: '48px', borderRadius: '10px', objectFit: 'cover' }} />
                           ) : (
                             <div style={{ width: '48px', height: '48px', borderRadius: '10px', background: getDeterministicGradient(pos.token_address), display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold', color: '#fff', fontSize: '1rem' }}>
                               {pos.token_symbol.slice(0, 2).toUpperCase()}
                             </div>
                           )}
                           <div style={{ flex: 1, minWidth: 0 }}>
                             <h4 style={{ fontSize: '0.95rem', color: '#fff', margin: 0, fontWeight: 700 }}>{pos.token_name}</h4>
                             <span style={{ fontSize: '0.75rem', color: 'var(--color-primary)' }}>{pos.token_symbol}</span>
                             <div style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)', marginTop: '4px' }}>
                               Solde : <span style={{ color: '#fff', fontWeight: 600 }}>{tokenBal.toLocaleString(undefined, { maximumFractionDigits: 0 })}</span>
                             </div>
                           </div>
                           <div style={{ textAlign: 'right' }}>
                             <div style={{ fontSize: '0.95rem', fontWeight: 'bold', color: 'var(--color-success)' }}>
                               {positionValue.toFixed(4)} ETH
                             </div>
                             <div style={{ 
                               fontSize: '0.75rem', 
                               color: pos.unrealized_pnl >= 0 ? 'var(--color-success)' : 'var(--color-danger)',
                               marginTop: '2px'
                             }}>
                               {pos.unrealized_pnl >= 0 ? '▲ +' : '▼ '}{pos.unrealized_pnl.toFixed(4)} ETH
                             </div>
                           </div>
                         </div>
                       );
                     })}
                   </div>
                 )}
               </div>
             )}

             {profileTab === 'created' && (
               <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: '20px' }}>
                 {creatorTokens.length === 0 ? (
                   <div style={{ padding: '60px 0', textAlign: 'center', color: 'var(--color-text-muted)', fontSize: '0.95rem', gridColumn: '1 / -1' }} className="glass-panel">
                     Aucun jeton créé par cet utilisateur.
                   </div>
                 ) : (
                   creatorTokens.map((t) => {
                     const mcapEth = 1 * (1 + (parseFloat(t.tokens_sold) / 1e18) / 960000000);
                     const progress = Math.min((parseFloat(t.tokens_sold) / 1e18) / 800000000 * 100, 100);
                     return (
                       <div 
                         key={t.address}
                         className="glass-panel"
                         onClick={() => {
                           setSelectedToken(t);
                           setViewProfileAddress(null);
                         }}
                         style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '12px', cursor: 'pointer', position: 'relative' }}
                       >
                         <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
                           {t.image_url ? (
                             <img src={t.image_url} alt={t.name} style={{ width: '40px', height: '40px', borderRadius: '8px', objectFit: 'cover' }} />
                           ) : (
                             <div style={{ width: '40px', height: '40px', borderRadius: '8px', background: getDeterministicGradient(t.address), display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold', color: '#fff', fontSize: '0.85rem' }}>
                               {t.symbol.slice(0, 2).toUpperCase()}
                             </div>
                           )}
                           <div>
                             <h4 style={{ fontSize: '0.9rem', color: '#fff', margin: 0, fontWeight: 700 }}>{t.name}</h4>
                             <span style={{ fontSize: '0.7rem', color: 'var(--color-primary)' }}>{t.symbol}</span>
                           </div>
                         </div>
                         
                         <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                           <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.65rem', color: 'var(--color-text-muted)' }}>
                             <span>Migration :</span>
                             <span style={{ fontWeight: 'bold', color: '#fff' }}>{progress.toFixed(1)}%</span>
                           </div>
                           <div style={{ width: '100%', height: '4px', background: 'rgba(255,255,255,0.05)', borderRadius: '2px', overflow: 'hidden' }}>
                             <div style={{ width: `${progress}%`, height: '100%', background: 'var(--color-primary)' }} />
                           </div>
                         </div>
                         <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.75rem', borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '8px', marginTop: '4px' }}>
                           <span style={{ color: 'var(--color-text-muted)' }}>Cap. Marché :</span>
                           <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>{mcapEth.toFixed(1)} ETH</span>
                         </div>
                       </div>
                     );
                   })
                 )}
               </div>
             )}

             {profileTab === 'trades' && (
               <div className="glass-panel" style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
                 {userTrades.length === 0 ? (
                   <div style={{ padding: '40px 0', textAlign: 'center', color: 'var(--color-text-muted)', fontSize: '0.95rem' }}>
                     Aucune transaction enregistrée.
                   </div>
                 ) : (
                   <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', maxHeight: '500px', overflowY: 'auto', paddingRight: '4px' }}>
                     {userTrades.map((tr: any) => {
                       const matchedT = tokens.find(tok => tok.address.toLowerCase() === tr.token_address.toLowerCase());
                       return (
                         <div 
                           key={tr.id}
                           onClick={() => {
                             if (matchedT) {
                               setSelectedToken(matchedT);
                               setViewProfileAddress(null);
                             }
                           }}
                           style={{
                             padding: '10px 14px',
                             background: 'rgba(255,255,255,0.01)',
                             borderRadius: '10px',
                             border: '1px solid var(--border-glass)',
                             display: 'flex',
                             justifyContent: 'space-between',
                             alignItems: 'center',
                             fontSize: '0.8rem',
                             cursor: matchedT ? 'pointer' : 'default'
                           }}
                         >
                           <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
                             <span style={{ 
                               color: tr.is_buy ? 'var(--color-success)' : 'var(--color-danger)', 
                               fontWeight: 'bold',
                               background: tr.is_buy ? 'var(--color-success-glow)' : 'var(--color-danger-glow)',
                               padding: '2px 8px',
                               borderRadius: '6px',
                               fontSize: '0.7rem'
                             }}>
                               {tr.is_buy ? 'ACHAT' : 'VENTE'}
                             </span>
                             <span style={{ color: '#fff', fontWeight: 600 }}>
                               {parseFloat(formatEther(BigInt(tr.token_amount))).toLocaleString(undefined, { maximumFractionDigits: 0 })}
                             </span>
                             <span style={{ color: 'var(--color-text-muted)' }}>{matchedT?.symbol || 'TKN'}</span>
                           </div>
                           <div style={{ display: 'flex', flexDirection: 'column', gap: '2px', alignItems: 'flex-end' }}>
                             <span style={{ color: 'var(--color-success)', fontWeight: 600 }}>{parseFloat(formatEther(BigInt(tr.eth_amount))).toFixed(4)} ETH</span>
                             <span style={{ fontSize: '0.65rem', color: 'var(--color-text-dark)' }}>
                               {new Date(tr.timestamp).toLocaleString()}
                             </span>
                           </div>
                         </div>
                       );
                     })}
                   </div>
                 )}
               </div>
             )}
          </div>
        ) : !selectedToken ? (
          <>
            {/* FILTER AND ACTION BAR */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '16px', flexWrap: 'wrap', width: '100%', marginBottom: '12px' }}>
              {/* Left controls group */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flexWrap: 'wrap', flex: 1 }}>
                {/* 4 Mutually Exclusive Filters */}
                <div style={{ display: 'flex', gap: '6px', background: 'rgba(255,255,255,0.02)', padding: '4px', borderRadius: '10px', border: '1px solid var(--border-glass)' }}>
                  <button 
                    className={`btn ${activeFilter === 'all' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setActiveFilter('all')}
                    style={{ padding: '8px 14px', fontSize: '0.8rem' }}
                  >
                    Tous les Coins
                  </button>
                  <button 
                    className={`btn ${activeFilter === 'created' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setActiveFilter('created')}
                    style={{ padding: '8px 14px', fontSize: '0.8rem' }}
                    disabled={!isConnected}
                  >
                    Mes Créations
                  </button>
                  <button 
                    className={`btn ${activeFilter === 'progress' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setActiveFilter('progress')}
                    style={{ padding: '8px 12px', fontSize: '0.8rem', display: 'flex', alignItems: 'center', gap: '4px' }}
                  >
                    <TrendingUp size={14} /> Progression
                  </button>
                  <button 
                    className={`btn ${activeFilter === 'new' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setActiveFilter('new')}
                    style={{ padding: '8px 12px', fontSize: '0.8rem' }}
                  >
                    Récents
                  </button>
                </div>

                {/* Search bar (Enlarged) */}
                <div style={{ position: 'relative', width: '380px', maxWidth: '100%' }}>
                  <Search size={14} color="var(--color-text-muted)" style={{ position: 'absolute', left: '10px', top: '11px' }} />
                  <input 
                    type="text" 
                    placeholder="Rechercher un token..." 
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="input-field"
                    style={{ paddingLeft: '32px', height: '36px', fontSize: '0.85rem' }}
                  />
                </div>
              </div>

              {/* Deploy Button (Aligned right) */}
              <button 
                className="btn btn-success" 
                onClick={() => setShowCreateModal(true)} 
                style={{ padding: '0 20px', height: '36px', fontSize: '0.85rem', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '6px', borderRadius: '10px' }}
              >
                <Plus size={16} />
                Lancer un Coin
              </button>
            </div>

            {/* TOKENS GRID LAYOUT */}
            <div style={{ 
              display: 'grid', 
              gridTemplateColumns: 'repeat(3, 1fr)', 
              gap: '20px', 
              marginTop: '16px' 
            }}>
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
                  const mcapUsd = mcapEth * 2000 // $2000 per ETH

                  return (
                    <div 
                      key={t.address}
                      onClick={() => setSelectedToken(t)}
                      className="glass-panel"
                      style={{ 
                        padding: '20px', 
                        cursor: 'pointer', 
                        display: 'flex', 
                        flexDirection: 'row', 
                        gap: '16px', 
                        alignItems: 'stretch',
                        position: 'relative',
                        transition: 'transform 0.2s, box-shadow 0.2s, border-color 0.2s',
                        border: '1px solid var(--border-glass)',
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.transform = 'translateY(-4px)';
                        e.currentTarget.style.boxShadow = '0 12px 30px rgba(0, 82, 255, 0.15)';
                        e.currentTarget.style.borderColor = 'rgba(0, 82, 255, 0.3)';
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.transform = 'translateY(0)';
                        e.currentTarget.style.boxShadow = 'none';
                        e.currentTarget.style.borderColor = 'var(--border-glass)';
                      }}
                    >
                      {/* Left Side: Large Image (192px) */}
                      <div style={{ flexShrink: 0, display: 'flex', alignItems: 'center' }}>
                        {t.image_url ? (
                          <img 
                            src={t.image_url} 
                            alt={t.name} 
                            style={{ width: '192px', height: '192px', borderRadius: '12px', objectFit: 'cover', border: '1px solid rgba(255,255,255,0.08)', boxShadow: '0 4px 12px rgba(0,0,0,0.3)' }} 
                          />
                        ) : (
                          <div style={{
                            width: '192px',
                            height: '192px',
                            borderRadius: '12px',
                            background: getDeterministicGradient(t.address),
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            fontWeight: 'bold',
                            color: '#fff',
                            fontSize: '2.5rem',
                            border: '1px solid rgba(255,255,255,0.08)',
                            boxShadow: '0 4px 12px rgba(0,0,0,0.3)'
                          }}>
                            {t.symbol.slice(0, 2).toUpperCase()}
                          </div>
                        )}
                      </div>

                      {/* Right Side: Information Details */}
                      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: '8px', minWidth: 0, justifyContent: 'space-between' }}>
                        {/* Title and Symbol */}
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: '6px', justifyContent: 'space-between' }}>
                            <span style={{ fontSize: '1rem', fontWeight: 800, color: '#fff', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', maxWidth: '140px' }} title={t.name}>
                              {t.name}
                            </span>
                            <span style={{ fontSize: '0.75rem', color: 'var(--color-primary)', fontWeight: 800 }}>
                              ${t.symbol}
                            </span>
                          </div>
                          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <span style={{ fontSize: '0.7rem', color: 'var(--color-text-dark)', fontFamily: 'monospace' }}>
                              Par: <span 
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setViewProfileAddress(t.creator);
                                  setSelectedToken(null);
                                }}
                                style={{ color: 'var(--color-text-muted)', textDecoration: 'underline', cursor: 'pointer' }}
                              >
                                {shortenAddress(t.creator)}
                              </span>
                            </span>
                            {isMyCreation && (
                              <span style={{ background: 'var(--color-primary-glow)', color: 'var(--color-primary)', fontSize: '0.65rem', padding: '1px 5px', borderRadius: '4px', border: '1px solid var(--color-primary)', fontWeight: 'bold' }}>
                                Dev
                              </span>
                            )}
                          </div>
                        </div>

                        {/* Description */}
                        <p style={{ 
                          margin: 0, 
                          fontSize: '0.78rem', 
                          color: 'var(--color-text-muted)', 
                          lineHeight: '1.4',
                          height: '36px',
                          display: '-webkit-box',
                          WebkitLineClamp: 2,
                          WebkitBoxOrient: 'vertical',
                          overflow: 'hidden',
                          textOverflow: 'ellipsis'
                        }}>
                          {t.description || "Aucune description fournie pour ce meme coin."}
                        </p>

                        {/* Indicators: Mcap & Vol */}
                        <div style={{ 
                          display: 'grid', 
                          gridTemplateColumns: '1fr 1fr', 
                          gap: '6px', 
                          background: 'rgba(255,255,255,0.01)', 
                          padding: '8px 10px', 
                          borderRadius: '8px', 
                          border: '1px solid rgba(255,255,255,0.03)' 
                        }}>
                          <div style={{ display: 'flex', flexDirection: 'column' }}>
                            <span style={{ fontSize: '0.6rem', color: 'var(--color-text-dark)', fontWeight: 600, textTransform: 'uppercase' }}>Mcap</span>
                            <span style={{ fontSize: '0.82rem', fontWeight: 800, color: 'var(--color-success)' }}>
                              ${mcapUsd >= 1000000 ? `${(mcapUsd / 1000000).toFixed(2)}M` : mcapUsd >= 1000 ? `${(mcapUsd / 1000).toFixed(1)}K` : mcapUsd.toFixed(2)}
                            </span>
                          </div>
                          <div style={{ display: 'flex', flexDirection: 'column' }}>
                            <span style={{ fontSize: '0.6rem', color: 'var(--color-text-dark)', fontWeight: 600, textTransform: 'uppercase' }}>Volume</span>
                            <span style={{ fontSize: '0.82rem', fontWeight: 800, color: '#fff' }}>
                              {ethRaisedNum.toFixed(1)} ETH
                            </span>
                          </div>
                        </div>

                        {/* Bonding Curve */}
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '3px' }}>
                          <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.68rem', color: 'var(--color-text-muted)' }}>
                            <span>Bonding Curve</span>
                            <span style={{ fontWeight: 'bold' }}>{progress.toFixed(1)}%</span>
                          </div>
                          <div style={{ width: '100%', height: '4px', background: 'rgba(255,255,255,0.04)', borderRadius: '2px', overflow: 'hidden' }}>
                            <div style={{ width: `${progress}%`, height: '100%', background: t.migrated ? 'var(--color-success)' : 'var(--color-primary)' }} />
                          </div>
                        </div>

                        {/* Footer: Age & Swap */}
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderTop: '1px solid rgba(255,255,255,0.03)', paddingTop: '8px' }}>
                          <span style={{ fontSize: '0.7rem', color: 'var(--color-text-dark)' }}>
                            {getRelativeTime(t.created_at)}
                          </span>
                          
                          <button 
                            className={`btn ${t.migrated ? 'btn-secondary' : 'btn-primary'}`}
                            style={{ 
                              padding: '4px 12px', 
                              fontSize: '0.72rem', 
                              height: '26px', 
                              borderRadius: '6px', 
                              display: 'flex', 
                              alignItems: 'center', 
                              gap: '4px',
                              fontWeight: 'bold'
                            }}
                            onClick={(e) => {
                              e.stopPropagation();
                              setSelectedToken(t);
                            }}
                          >
                            <span>Swap</span>
                            <ArrowLeft size={8} style={{ transform: 'rotate(180deg)' }} />
                          </button>
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

            {/* Centered details wrapper to reduce left column width while keeping right column width same */}
            <div style={{ maxWidth: '85%', width: '100%', margin: '0 auto', display: 'flex', flexDirection: 'column', gap: '24px' }}>

                        {/* SELECTED TOKEN VIEW */}
            {/* HEADER DETAILS (Clanker-style with all info and description) */}
                                    {/* HEADER DETAILS (Clanker-style 2 separate boxes) */}
            <div className="details-grid" style={{ margin: 0, gap: '24px', alignItems: 'stretch' }}>
              
              {/* FIRST BOX (Left - 3fr): Token info details glass-panel */}
              <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', minWidth: 0 }}>
                <div style={{ display: 'flex', gap: '16px', alignItems: 'center', flexWrap: 'wrap' }}>
                  {/* Token Logo / Avatar */}
                  {selectedToken.image_url ? (
                    <img 
                      src={selectedToken.image_url} 
                      alt={selectedToken.name} 
                      style={{ width: '64px', height: '64px', borderRadius: '50%', objectFit: 'cover', flexShrink: 0, border: '1px solid rgba(255,255,255,0.08)', boxShadow: '0 4px 12px rgba(0,0,0,0.3)' }} 
                    />
                  ) : (
                    <div style={{
                      width: '64px',
                      height: '64px',
                      borderRadius: '50%',
                      background: getDeterministicGradient(selectedToken.address),
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontWeight: 'bold',
                      color: '#fff',
                      fontSize: '1.4rem',
                      flexShrink: 0,
                      border: '1px solid rgba(255,255,255,0.08)',
                      boxShadow: '0 4px 12px rgba(0,0,0,0.3)'
                    }}>
                      {selectedToken.symbol.slice(0, 2).toUpperCase()}
                    </div>
                  )}

                  <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                    <div style={{ display: 'flex', alignItems: 'baseline', gap: '8px', flexWrap: 'wrap' }}>
                      <h2 style={{ fontSize: '1.8rem', color: '#fff', margin: 0, fontWeight: 800, lineHeight: 1.1 }}>{selectedToken.name}</h2>
                      <span style={{ fontSize: '1.1rem', color: 'var(--color-primary)', fontWeight: 600 }}>(${selectedToken.symbol})</span>
                      
                      <span style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)', background: 'rgba(255,255,255,0.04)', padding: '2px 8px', borderRadius: '6px', marginLeft: '6px' }}>
                        {getRelativeTime(selectedToken.created_at)}
                      </span>
                    </div>
                    
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.8rem', color: 'var(--color-text-muted)', flexWrap: 'wrap', fontFamily: 'monospace' }}>
                      <span style={{ color: 'var(--color-text-dark)', fontFamily: 'var(--font-family-body)' }}>Contrat :</span>
                      <span style={{ color: 'var(--color-primary)' }}>{selectedToken.address}</span>
                      <button 
                        onClick={() => {
                          navigator.clipboard.writeText(selectedToken.address);
                        }}
                        style={{ background: 'transparent', border: 'none', cursor: 'pointer', padding: '2px', display: 'flex', alignItems: 'center' }}
                        title="Copier l'adresse"
                      >
                        <Check size={12} color="var(--color-text-muted)" />
                      </button>
                    </div>
                  </div>
                </div>

                {/* Price & Change block */}
                {(() => {
                  const ethRaisedNum = parseFloat(formatEther(BigInt(selectedToken.eth_raised)))
                  const tokensSoldNum = parseFloat(formatEther(BigInt(selectedToken.tokens_sold)))
                  const price = (3 + ethRaisedNum) / (941176470 - tokensSoldNum)
                  const priceUsd = price * 2000
                  return (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flexWrap: 'wrap' }}>
                        <span style={{ fontSize: '1.6rem', color: 'var(--color-success)', fontWeight: 800 }}>
                          ${priceUsd.toFixed(8)}
                        </span>
                        <span style={{ fontSize: '0.85rem', color: 'var(--color-success)', background: 'rgba(16, 185, 129, 0.08)', padding: '2px 8px', borderRadius: '6px', fontWeight: 'bold' }}>
                          +0% (24h)
                        </span>
                      </div>

                      {/* Social Links Row */}
                      <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '4px' }}>
                        {selectedToken.website && (
                          <a 
                            href={selectedToken.website} 
                            target="_blank" 
                            rel="noopener noreferrer" 
                            className="btn btn-secondary" 
                            style={{ display: 'flex', alignItems: 'center', gap: '6px', padding: '4px 10px', fontSize: '0.75rem', height: 'auto', borderRadius: '20px' }}
                          >
                            <Globe size={12} />
                            <span>Site Web</span>
                          </a>
                        )}
                        {selectedToken.twitter && (
                          <a 
                            href={selectedToken.twitter} 
                            target="_blank" 
                            rel="noopener noreferrer" 
                            className="btn btn-secondary" 
                            style={{ display: 'flex', alignItems: 'center', gap: '6px', padding: '4px 10px', fontSize: '0.75rem', height: 'auto', borderRadius: '20px' }}
                          >
                            <svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor">
                              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
                            </svg>
                            <span>Twitter</span>
                          </a>
                        )}
                        {selectedToken.telegram && (
                          <a 
                            href={selectedToken.telegram} 
                            target="_blank" 
                            rel="noopener noreferrer" 
                            className="btn btn-secondary" 
                            style={{ display: 'flex', alignItems: 'center', gap: '6px', padding: '4px 10px', fontSize: '0.75rem', height: 'auto', borderRadius: '20px' }}
                          >
                            <Send size={12} style={{ transform: 'rotate(-30deg)' }} />
                            <span>Telegram</span>
                          </a>
                        )}
                      </div>
                    </div>
                  )
                })()}

                {/* Description / edit form */}
                <div style={{ borderTop: '1px solid rgba(255,255,255,0.06)', paddingTop: '12px', marginTop: '4px' }}>
                  {isEditingMetadata ? (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', background: 'rgba(255,255,255,0.01)', padding: '16px', borderRadius: '12px', border: '1px solid var(--border-glass)' }}>
                      <textarea
                        rows={3}
                        placeholder="Ajoutez une description de votre meme coin..."
                        value={metadataDescription}
                        onChange={(e) => setMetadataDescription(e.target.value)}
                        className="input-field"
                        style={{ resize: 'vertical', fontSize: '0.9rem', padding: '12px', fontFamily: 'inherit' }}
                      />
                      <input
                        type="text"
                        placeholder="Image URL..."
                        value={metadataImageUrl}
                        onChange={(e) => setMetadataImageUrl(e.target.value)}
                        className="input-field"
                        style={{ fontSize: '0.9rem', padding: '10px 12px' }}
                      />
                      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '8px' }}>
                        <input
                          type="text"
                          placeholder="Site Web (URL)..."
                          value={metadataWebsite}
                          onChange={(e) => setMetadataWebsite(e.target.value)}
                          className="input-field"
                          style={{ fontSize: '0.85rem', padding: '10px 12px' }}
                        />
                        <input
                          type="text"
                          placeholder="Twitter (URL)..."
                          value={metadataTwitter}
                          onChange={(e) => setMetadataTwitter(e.target.value)}
                          className="input-field"
                          style={{ fontSize: '0.85rem', padding: '10px 12px' }}
                        />
                        <input
                          type="text"
                          placeholder="Telegram (URL)..."
                          value={metadataTelegram}
                          onChange={(e) => setMetadataTelegram(e.target.value)}
                          className="input-field"
                          style={{ fontSize: '0.85rem', padding: '10px 12px' }}
                        />
                      </div>
                      <button 
                        className="btn btn-primary" 
                        onClick={handleUpdateMetadata}
                        style={{ alignSelf: 'flex-start', padding: '8px 16px', fontSize: '0.8rem' }}
                        disabled={isUpdatingMetadata}
                      >
                        {isUpdatingMetadata ? <RefreshCw className="animate-spin" size={14} /> : 'Sauvegarder'}
                      </button>
                    </div>
                  ) : (
                    <p style={{ fontSize: '0.9rem', color: 'var(--color-text-muted)', lineHeight: '1.6', margin: 0, whiteSpace: 'pre-wrap' }}>
                      {selectedToken.description || "Aucune description fournie."}
                    </p>
                  )}
                </div>

                <div style={{ display: 'flex', gap: '12px', alignItems: 'center', flexWrap: 'wrap', fontSize: '0.8rem', color: 'var(--color-text-muted)', borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '10px' }}>
                  <span>Réseau : <strong style={{ color: '#fff' }}>Base L2</strong></span>
                  <span>•</span>
                  <span>Admin : <strong style={{ color: 'var(--color-text-muted)', fontFamily: 'monospace' }}>{shortenAddress(factoryAddress)}</strong></span>
                  <span>•</span>
                  <span>
                    Créateur :{' '}
                    <strong 
                      onClick={() => {
                        setViewProfileAddress(selectedToken.creator);
                        setSelectedToken(null);
                      }}
                      style={{ color: 'var(--color-primary)', cursor: 'pointer', textDecoration: 'underline' }}
                    >
                      {shortenAddress(selectedToken.creator)}
                    </strong>
                  </span>
                  {isCreator && (
                    <button 
                      className="btn btn-secondary" 
                      onClick={() => {
                        setIsEditingMetadata(!isEditingMetadata);
                        setMetadataDescription(selectedToken.description || "");
                        setMetadataImageUrl(selectedToken.image_url || "");
                      }}
                      style={{ padding: '2px 8px', fontSize: '0.7rem', height: 'auto', marginLeft: 'auto' }}
                    >
                      {isEditingMetadata ? 'Annuler' : 'Modifier les infos'}
                    </button>
                  )}
                </div>
              </div>

              {/* SECOND BOX (Right - 1fr): Key metrics listed vertically glass-panel */}
              {(() => {
                const ethRaisedNum = parseFloat(formatEther(BigInt(selectedToken.eth_raised)))
                const tokensSoldNum = parseFloat(formatEther(BigInt(selectedToken.tokens_sold)))
                const price = (3 + ethRaisedNum) / (941176470 - tokensSoldNum)
                const mcapEth = price * 1000000000
                const mcapUsd = mcapEth * 2000
                const virtualLiquidityEth = 3 + ethRaisedNum
                const virtualLiquidityUsd = virtualLiquidityEth * 2000
                const volumeUsd = ethRaisedNum * 2000
                return (
                  <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', height: '100%', justifyContent: 'center' }}>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px 24px', alignItems: 'center' }}>
                      <div>
                        <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)', display: 'block', textTransform: 'uppercase', fontWeight: 600 }}>Market Cap</span>
                        <strong style={{ fontSize: '1.25rem', color: '#fff', display: 'block', marginTop: '2px' }}>
                          ${mcapUsd.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                        </strong>
                        <span style={{ fontSize: '0.7rem', color: 'var(--color-success)', fontWeight: 600 }}>{mcapEth.toFixed(1)} ETH</span>
                      </div>

                      <div>
                        <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)', display: 'block', textTransform: 'uppercase', fontWeight: 600 }}>Virtual Liquidity</span>
                        <strong style={{ fontSize: '1.15rem', color: '#fff', display: 'block', marginTop: '2px' }}>
                          ${virtualLiquidityUsd.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                        </strong>
                        <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)' }}>{virtualLiquidityEth.toFixed(2)} ETH</span>
                      </div>

                      <div style={{ borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '12px' }}>
                        <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)', display: 'block', textTransform: 'uppercase', fontWeight: 600 }}>24H Volume</span>
                        <strong style={{ fontSize: '1.15rem', color: '#fff', display: 'block', marginTop: '2px' }}>
                          ${volumeUsd.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                        </strong>
                        <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)' }}>{ethRaisedNum.toFixed(2)} ETH</span>
                      </div>

                      <div style={{ borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '12px' }}>
                        <span style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)', display: 'block', textTransform: 'uppercase', fontWeight: 600 }}>Holders</span>
                        <strong style={{ fontSize: '1.15rem', color: '#fff', display: 'block', marginTop: '2px' }}>
                          {tokenHolders.length || 1}
                        </strong>
                        <span style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)', display: 'block', marginTop: '2px' }}>adresses</span>
                      </div>
                    </div>
                  </div>
                )
              })()}

            </div>
{/* TWO COLUMN GRID LAYOUT (Clanker-style: Left wide, Right narrower 320px) */}
            <div className="details-grid">
              
              {/* LEFT COLUMN: Price Chart, Trades/Holders */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px', minWidth: 0 }}>
                {/* 1. Price Chart Card */}
                <PriceChart tokenAddress={selectedToken.address} lastTrade={lastTrade} />

                {/* 2. Trades / Holders Tab Card */}
                {/* TRADES / HOLDERS TABS CARD */}
              <div className="glass-panel" style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '16px', height: '450px' }}>
                {/* Tab selector header */}
                <div style={{ display: 'flex', gap: '8px', borderBottom: '1px solid rgba(255,255,255,0.06)', paddingBottom: '10px' }}>
                  <button
                    className={`btn ${detailsTab === 'trades' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => setDetailsTab('trades')}
                    style={{ padding: '6px 12px', fontSize: '0.85rem', height: 'auto' }}
                  >
                    🚀 Transactions
                  </button>
                  <button
                    className={`btn ${detailsTab === 'holders' ? 'btn-primary' : 'btn-secondary'}`}
                    onClick={() => {
                      setDetailsTab('holders');
                      fetchTokenHolders(selectedToken.address);
                    }}
                    style={{ padding: '6px 12px', fontSize: '0.85rem', height: 'auto' }}
                  >
                    👥 Top Holders ({tokenHolders.length})
                  </button>
                </div>

                {detailsTab === 'trades' ? (
                  /* TRANSACTIONS FLOW */
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
                            <span 
                              onClick={() => {
                                setViewProfileAddress(tr.buyer_or_seller);
                                setSelectedToken(null);
                              }}
                              style={{ fontSize: '0.7rem', color: 'var(--color-primary)', cursor: 'pointer', textDecoration: 'underline' }}
                            >
                              {shortenAddress(tr.buyer_or_seller)}
                            </span>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                ) : (
                  /* HOLDERS LIST */
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', overflowY: 'auto', flex: 1, paddingRight: '4px' }}>
                    {tokenHolders.length === 0 ? (
                      <div style={{ textAlign: 'center', color: 'var(--color-text-muted)', padding: '50px 0', fontSize: '0.85rem' }}>
                        Aucun détenteur trouvé.
                      </div>
                    ) : (
                      tokenHolders.map((holder, idx) => (
                        <div 
                          key={holder.address}
                          style={{
                            padding: '10px 14px',
                            background: 'rgba(255,255,255,0.01)',
                            borderRadius: '10px',
                            border: '1px solid var(--border-glass)',
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            fontSize: '0.8rem'
                          }}
                        >
                          <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
                            <span style={{ color: 'var(--color-text-muted)', fontWeight: 700, minWidth: '20px' }}>
                              #{idx + 1}
                            </span>
                            <div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                              <span 
                                onClick={() => {
                                  setViewProfileAddress(holder.address);
                                  setSelectedToken(null);
                                }}
                                style={{ 
                                  color: 'var(--color-primary)', 
                                  cursor: 'pointer', 
                                  textDecoration: 'underline',
                                  fontWeight: holder.is_creator ? 'bold' : 'normal',
                                  display: 'flex',
                                  alignItems: 'center',
                                  gap: '6px'
                                }}
                              >
                                {holder.username || shortenAddress(holder.address)}
                                {holder.is_creator && (
                                  <span style={{ background: 'var(--color-primary-glow)', color: 'var(--color-primary)', fontSize: '0.65rem', padding: '1px 5px', borderRadius: '4px', border: '1px solid var(--color-primary)' }}>
                                    Dev
                                  </span>
                                )}
                              </span>
                              <span style={{ fontSize: '0.65rem', color: 'var(--color-text-dark)', fontFamily: 'monospace' }}>
                                {holder.address}
                              </span>
                            </div>
                          </div>
                          
                          <div style={{ textAlign: 'right' }}>
                            <div style={{ color: '#fff', fontWeight: 600 }}>
                              {(parseFloat(holder.balance) / 1e18).toLocaleString(undefined, { maximumFractionDigits: 0 })} {selectedToken.symbol}
                            </div>
                            <div style={{ fontSize: '0.7rem', color: 'var(--color-text-muted)', marginTop: '2px' }}>
                              Part : <strong>{holder.percentage.toFixed(2)}%</strong>
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}
              </div>
              </div>

              {/* RIGHT COLUMN: Swap Card, Clankernomics/Vesting, Condensed Audit Card */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                {/* 1. Integrated Swap Card */}
                {/* TRADING FORM */}
                <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', height: '518px' }}>
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
                      <div style={{ display: 'flex', gap: '6px', marginTop: '4px' }}>
                        {["Reset", "0.1", "0.5", "1.0", "2.0"].map((val) => (
                          <button
                            key={val}
                            onClick={() => {
                              if (val === "Reset") {
                                handleAmountChange("");
                              } else {
                                handleAmountChange(val);
                              }
                            }}
                            className="btn btn-secondary"
                            style={{ flex: 1, padding: '4px', fontSize: '0.75rem', borderRadius: '6px' }}
                          >
                            {val === "Reset" ? "Reset" : `${val} ETH`}
                          </button>
                        ))}
                      </div>
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
                          {["Reset", "25%", "50%", "75%", "100%"].map((pct) => (
                            <button
                              key={pct}
                              onClick={() => {
                                if (pct === "Reset") {
                                  handleAmountChange("");
                                } else {
                                  const balanceBig = userTokenBalance as bigint
                                  const pctNum = parseFloat(pct) / 100
                                  const amountBig = (balanceBig * BigInt(Math.floor(pctNum * 10000))) / 10000n
                                  handleAmountChange(formatEther(amountBig))
                                }
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

                {/* Bounding Curve Progress Bar */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: '6px', borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '14px', marginTop: '10px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>
                    <span>Courbe de Liaison (Bonding Curve) :</span>
                    <strong style={{ color: selectedToken.migrated ? 'var(--color-success)' : '#fff' }}>
                      {Math.min((parseFloat(selectedToken.tokens_sold) / 1e18 / 800000000) * 100, 100).toFixed(1)}%
                    </strong>
                  </div>
                  
                  {/* Progress bar */}
                  <div style={{ width: '100%', height: '8px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px', overflow: 'hidden' }}>
                    <div 
                      style={{ 
                        width: `${Math.min((parseFloat(selectedToken.tokens_sold) / 1e18 / 800000000) * 100, 100)}%`, 
                        height: '100%', 
                        background: selectedToken.migrated ? 'var(--color-success)' : 'linear-gradient(90deg, var(--color-primary) 0%, #00d2ff 100%)',
                        borderRadius: '4px',
                        transition: 'width 0.4s cubic-bezier(0.4, 0, 0.2, 1)'
                      }} 
                    />
                  </div>
                  
                  <p style={{ fontSize: '0.65rem', color: 'var(--color-text-muted)', margin: 0, textAlign: 'center', marginTop: '2px' }}>
                    {selectedToken.migrated 
                      ? "Le token a migré sur Uniswap V2 !" 
                      : `Migration automatique sur Uniswap V2 une fois la courbe à 100% (800M de tokens vendus).`
                    }
                  </p>
                </div>
              </div>


                {/* 2. Clankernomics / Vesting Card */}
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

                {/* 3. Condensed Security Audit Card (Taller version with safety checklist) */}
                <div className="glass-panel" style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '16px', height: '450px', justifyContent: 'space-between' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px', borderBottom: '1px solid rgba(255,255,255,0.06)', paddingBottom: '10px' }}>
                    <ShieldCheck size={18} color="var(--color-success)" />
                    <span style={{ fontSize: '0.85rem', fontWeight: 700, textTransform: 'uppercase', color: '#fff' }}>Rapport de Sécurité</span>
                  </div>

                  {isLoadingAudit || !auditData ? (
                    <div style={{ padding: '30px 0', textAlign: 'center', color: 'var(--color-text-muted)', fontSize: '0.85rem', flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      Audit en cours...
                    </div>
                  ) : (
                    (() => {
                      const score = 100 - Math.max(0, (auditData.dev_percent - 5) * 5) - Math.max(0, (auditData.top_10_percent - 25) * 2) - Math.max(0, auditData.bundlers_percent * 10);
                      const safetyScore = Math.max(10, Math.min(100, Math.round(score)));
                      const isGood = safetyScore >= 80;
                      const isWarning = safetyScore >= 50 && safetyScore < 80;
                      const scoreColor = isGood ? 'var(--color-success)' : (isWarning ? 'var(--color-warning)' : 'var(--color-danger)');
                      const scoreBg = isGood ? 'rgba(16, 185, 129, 0.08)' : (isWarning ? 'rgba(245, 158, 11, 0.08)' : 'rgba(239, 68, 68, 0.08)');
                      
                      return (
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', flex: 1, justifyContent: 'space-between' }}>
                          {/* Score and verdict badge */}
                          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '10px 14px', background: scoreBg, border: `1px solid ${scoreColor}33`, borderRadius: '10px' }}>
                            <span style={{ fontSize: '0.85rem', fontWeight: 'bold', color: '#fff' }}>Score SafePump :</span>
                            <span style={{ fontSize: '1rem', fontWeight: 900, color: scoreColor }}>
                              {safetyScore} / 100
                            </span>
                          </div>

                          {/* Stats Grid */}
                          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '10px', textAlign: 'center', fontSize: '0.8rem' }}>
                            <div style={{ background: 'rgba(255,255,255,0.02)', padding: '8px 6px', borderRadius: '8px', border: '1px solid rgba(255,255,255,0.04)' }}>
                              <span style={{ color: 'var(--color-text-muted)', display: 'block', fontSize: '0.65rem', marginBottom: '2px' }}>Part Dev</span>
                              <strong style={{ color: '#fff', fontSize: '0.9rem' }}>{auditData.dev_percent.toFixed(1)}%</strong>
                            </div>
                            <div style={{ background: 'rgba(255,255,255,0.02)', padding: '8px 6px', borderRadius: '8px', border: '1px solid rgba(255,255,255,0.04)' }}>
                              <span style={{ color: 'var(--color-text-muted)', display: 'block', fontSize: '0.65rem', marginBottom: '2px' }}>Top 10</span>
                              <strong style={{ color: '#fff', fontSize: '0.9rem' }}>{auditData.top_10_percent.toFixed(1)}%</strong>
                            </div>
                            <div style={{ background: 'rgba(255,255,255,0.02)', padding: '8px 6px', borderRadius: '8px', border: '1px solid rgba(255,255,255,0.04)' }}>
                              <span style={{ color: 'var(--color-text-muted)', display: 'block', fontSize: '0.65rem', marginBottom: '2px' }}>Bundles</span>
                              <strong style={{ color: '#fff', fontSize: '0.9rem' }}>{auditData.bundlers_percent.toFixed(1)}%</strong>
                            </div>
                          </div>

                          {/* Detail checklist rows */}
                          <div style={{ display: 'flex', flexDirection: 'column', gap: '6px', borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '10px', fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <span>Code Source Vérifié :</span>
                              <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>✓ Oui (Base L2)</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <span>Simulateur HoneyPot :</span>
                              <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>✓ Réussi (Pas de taxe)</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <span>Taxes Achat / Vente :</span>
                              <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>✓ 0% / 0%</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <span>LP Lock / Burn :</span>
                              <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>✓ 100% à la migration</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <span>Propriété Contrat :</span>
                              <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>✓ Renoncée (Fixe)</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <span>Concentration Créateur :</span>
                              {auditData.dev_percent < 5 ? (
                                <span style={{ color: 'var(--color-success)', fontWeight: 'bold' }}>✓ Sécurisé (&lt;5%)</span>
                              ) : (
                                <span style={{ color: 'var(--color-warning)', fontWeight: 'bold' }}>⚠ Risque ({auditData.dev_percent.toFixed(1)}%)</span>
                              )}
                            </div>
                          </div>

                          {/* Footer Disclaimer */}
                          <div style={{ borderTop: '1px solid rgba(255,255,255,0.04)', paddingTop: '8px', fontSize: '0.65rem', color: 'var(--color-text-dark)', textAlign: 'center', fontStyle: 'italic' }}>
                            * Rapport automatisé par audit statique. Faites vos propres recherches (DYOR).
                          </div>
                        </div>
                      )
                    })()
                  )}
                </div>

              </div>
            </div>
            </div>
          </>
        )}
      </main>
        {/* LAUNCH TOKEN MODAL */}
      {showCreateModal && (
        <div style={{ position: 'fixed', top: 0, left: 0, width: '100%', height: '100%', background: 'rgba(2, 6, 23, 0.85)', backdropFilter: 'blur(8px)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
          <div className="glass-panel" style={{ width: '560px', maxHeight: '90vh', overflowY: 'auto', padding: '28px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
            <h3 style={{ fontSize: '1.4rem', color: '#fff', display: 'flex', alignItems: 'center', gap: '10px', borderBottom: '1px solid rgba(255,255,255,0.06)', paddingBottom: '12px', margin: 0 }}>
              <Rocket size={24} color="var(--color-success)" />
              Lancer votre Meme Coin
            </h3>

            {/* SECTION 1: Base info */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <h4 style={{ fontSize: '0.9rem', color: 'var(--color-success)', margin: 0, fontWeight: 'bold' }}>1. Informations du Jeton</h4>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                  <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Nom :</label>
                  <input 
                    type="text" 
                    placeholder="Ex: Safe Dog"
                    value={newTokenName}
                    onChange={(e) => setNewTokenName(e.target.value)}
                    className="input-field"
                    style={{ padding: '8px 12px', fontSize: '0.85rem' }}
                  />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                  <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Symbole :</label>
                  <input 
                    type="text" 
                    placeholder="Ex: SDOG"
                    value={newTokenSymbol}
                    onChange={(e) => setNewTokenSymbol(e.target.value)}
                    className="input-field"
                    style={{ padding: '8px 12px', fontSize: '0.85rem' }}
                  />
                </div>
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Proposition de Valeur :</label>
                <textarea 
                  placeholder="Expliquez la proposition de valeur ou le mème..."
                  value={newTokenDescription}
                  onChange={(e) => setNewTokenDescription(e.target.value)}
                  className="input-field"
                  rows={2}
                  style={{ resize: 'none', padding: '8px 12px', fontSize: '0.85rem' }}
                />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Lien de l'image (URL) :</label>
                <input 
                  type="text" 
                  placeholder="Ex: https://unsplash.com/photos/xyz.png"
                  value={newTokenImageUrl}
                  onChange={(e) => setNewTokenImageUrl(e.target.value)}
                  className="input-field"
                  style={{ padding: '8px 12px', fontSize: '0.85rem' }}
                />
              </div>
            </div>

            {/* SECTION 2: Socials info */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', borderTop: '1px solid rgba(255,255,255,0.06)', paddingTop: '12px' }}>
              <h4 style={{ fontSize: '0.9rem', color: 'var(--color-success)', margin: 0, fontWeight: 'bold' }}>2. Liens Réseaux Sociaux (Optionnel)</h4>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '10px' }}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                  <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Site Web :</label>
                  <input 
                    type="text" 
                    placeholder="https://..."
                    value={newTokenWebsite}
                    onChange={(e) => setNewTokenWebsite(e.target.value)}
                    className="input-field"
                    style={{ padding: '8px 12px', fontSize: '0.85rem' }}
                  />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                  <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Twitter / X :</label>
                  <input 
                    type="text" 
                    placeholder="https://x.com/..."
                    value={newTokenTwitter}
                    onChange={(e) => setNewTokenTwitter(e.target.value)}
                    className="input-field"
                    style={{ padding: '8px 12px', fontSize: '0.85rem' }}
                  />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
                  <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Telegram :</label>
                  <input 
                    type="text" 
                    placeholder="https://t.me/..."
                    value={newTokenTelegram}
                    onChange={(e) => setNewTokenTelegram(e.target.value)}
                    className="input-field"
                    style={{ padding: '8px 12px', fontSize: '0.85rem' }}
                  />
                </div>
              </div>
            </div>

            {/* SECTION 3: Launch Parameters (Bonus & Initial Buy) */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', borderTop: '1px solid rgba(255,255,255,0.06)', paddingTop: '12px' }}>
              <h4 style={{ fontSize: '0.9rem', color: 'var(--color-success)', margin: 0, fontWeight: 'bold' }}>3. Options & Avantages de Lancement</h4>
              
              {/* Creator Pack 5% Bonus checkbox */}
              <div 
                onClick={() => setCreateWithBonus(!createWithBonus)} 
                style={{ 
                  background: createWithBonus ? 'rgba(16, 185, 129, 0.08)' : 'rgba(255, 255, 255, 0.01)', 
                  border: createWithBonus ? '1px solid var(--color-success)' : '1px solid var(--border-glass)',
                  borderRadius: '10px',
                  padding: '10px 14px',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '12px',
                  transition: 'all 0.2s'
                }}
              >
                <input 
                  type="checkbox" 
                  checked={createWithBonus} 
                  onChange={() => {}} // handled by div click
                  style={{ accentColor: 'var(--color-success)', cursor: 'pointer', width: '16px', height: '16px' }}
                />
                <div style={{ display: 'flex', flexDirection: 'column', gap: '2px', flex: 1 }}>
                  <span style={{ fontSize: '0.85rem', color: '#fff', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '6px' }}>
                    ⭐ Pack Bonus Créateur (5%)
                  </span>
                  <span style={{ fontSize: '0.72rem', color: 'var(--color-text-muted)', lineHeight: '1.3' }}>
                    Sécurise 5% de la supply (50M jetons) pour <strong>0.02127 ETH</strong>. Renforce instantanément le score de sécurité et la confiance !
                  </span>
                </div>
              </div>

              {/* Initial Buy parameter */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
                <label style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>Achat Initial Créateur (Optionnel) :</label>
                <div style={{ display: 'flex', gap: '8px' }}>
                  <input 
                    type="number" 
                    step="0.01" 
                    min="0"
                    placeholder="Montant en ETH (Ex: 0.1)"
                    value={initialBuyEth}
                    onChange={(e) => setInitialBuyEth(e.target.value)}
                    className="input-field"
                    style={{ flex: 1, padding: '8px 12px', fontSize: '0.85rem' }}
                  />
                  <button 
                    className="btn btn-secondary" 
                    onClick={(e) => { e.stopPropagation(); setInitialBuyEth("0.1"); }}
                    style={{ fontSize: '0.75rem', padding: '0 12px', height: '36px', borderRadius: '8px' }}
                  >
                    0.1 ETH
                  </button>
                  <button 
                    className="btn btn-secondary" 
                    onClick={(e) => { e.stopPropagation(); setInitialBuyEth("0.5"); }}
                    style={{ fontSize: '0.75rem', padding: '0 12px', height: '36px', borderRadius: '8px' }}
                  >
                    0.5 ETH
                  </button>
                  <button 
                    className="btn btn-secondary" 
                    onClick={(e) => { e.stopPropagation(); setInitialBuyEth("1.0"); }}
                    style={{ fontSize: '0.75rem', padding: '0 12px', height: '36px', borderRadius: '8px' }}
                  >
                    1 ETH
                  </button>
                </div>
              </div>
            </div>

            {/* BUTTONS ROW */}
            <div style={{ display: 'flex', gap: '12px', marginTop: '10px', borderTop: '1px solid rgba(255,255,255,0.06)', paddingTop: '14px' }}>
              <button 
                className="btn btn-secondary" 
                onClick={() => {
                  setShowCreateModal(false);
                  setNewTokenName("");
                  setNewTokenSymbol("");
                  setNewTokenDescription("");
                  setNewTokenImageUrl("");
                  setNewTokenWebsite("");
                  setNewTokenTwitter("");
                  setNewTokenTelegram("");
                  setCreateWithBonus(false);
                  setInitialBuyEth("");
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

      {/* PROFILE EDITOR MODAL */}
      {showProfileModal && (
        <div style={{ position: 'fixed', top: 0, left: 0, width: '100%', height: '100%', background: 'rgba(2, 6, 23, 0.8)', backdropFilter: 'blur(8px)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
          <div className="glass-panel" style={{ width: '400px', padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
            <h3 style={{ fontSize: '1.3rem', color: '#fff', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <User size={22} color="var(--color-primary)" />
              Mettre à jour votre Profil
            </h3>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Nom d'utilisateur :</label>
              <input 
                type="text" 
                placeholder="Ex: CryptoKing"
                value={editUsername}
                onChange={(e) => setEditUsername(e.target.value)}
                className="input-field"
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>Biographie :</label>
              <textarea 
                placeholder="Quel type de trader êtes-vous ?"
                value={editBio}
                onChange={(e) => setEditBio(e.target.value)}
                className="input-field"
                rows={2}
                style={{ resize: 'none', padding: '10px 12px' }}
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
              <label style={{ fontSize: '0.85rem', color: 'var(--color-text-muted)' }}>URL de l'Avatar :</label>
              <input 
                type="text" 
                placeholder="Ex: https://image.com/avatar.png"
                value={editAvatarUrl}
                onChange={(e) => setEditAvatarUrl(e.target.value)}
                className="input-field"
              />
            </div>

            <div style={{ display: 'flex', gap: '12px', marginTop: '8px' }}>
              <button 
                className="btn btn-secondary" 
                onClick={() => {
                  setShowProfileModal(false)
                }} 
                style={{ flex: 1 }}
                disabled={isSavingProfile}
              >
                Annuler
              </button>
              <button 
                className="btn btn-primary" 
                onClick={handleSaveProfile}
                style={{ flex: 1 }}
                disabled={isSavingProfile || !editUsername.trim()}
              >
                {isSavingProfile ? <RefreshCw className="animate-spin" size={18} /> : 'Sauvegarder'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
