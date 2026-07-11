import { useEffect, useRef, useState } from 'react'
import { createChart, ColorType, CandlestickSeries } from 'lightweight-charts'
import type { ISeriesApi } from 'lightweight-charts'
import { formatEther } from 'viem'

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

interface PriceChartProps {
  tokenAddress: string
  lastTrade: Trade | null
}

interface CandleData {
  time: number
  open: number
  high: number
  low: number
  close: number
}

export default function PriceChart({ tokenAddress, lastTrade }: PriceChartProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const seriesRef = useRef<ISeriesApi<"Candlestick"> | null>(null)
  const [resolution, setResolution] = useState("1m")
  const [chartData, setChartData] = useState<CandleData[]>([])

  // Fetch historical candles
  const fetchCandles = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/v1/tokens/${tokenAddress}/candles?resolution=${resolution}`)
      const data = await res.json()
      if (Array.isArray(data)) {
        setChartData(data)
      } else {
        setChartData([])
      }
    } catch (e) {
      console.error("Failed to fetch candle data:", e)
      setChartData([])
    }
  }

  useEffect(() => {
    fetchCandles()
  }, [tokenAddress, resolution])

  // 1. Initialize the chart once on mount
  useEffect(() => {
    if (!chartContainerRef.current) return

    const handleResize = () => {
      chart.applyOptions({ width: chartContainerRef.current?.clientWidth })
    }

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
        textColor: '#94a3b8',
      },
      grid: {
        vertLines: { color: 'rgba(255, 255, 255, 0.04)' },
        horzLines: { color: 'rgba(255, 255, 255, 0.04)' },
      },
      width: chartContainerRef.current.clientWidth,
      height: 300,
      timeScale: {
        timeVisible: true,
        secondsVisible: false,
      },
    })

    const candlestickSeries = chart.addSeries(CandlestickSeries, {
      upColor: '#10b981',
      downColor: '#ef4444',
      borderVisible: false,
      wickUpColor: '#10b981',
      wickDownColor: '#ef4444',
      priceFormat: {
        type: 'price',
        precision: 10,
        minMove: 0.0000000001,
      },
    })

    seriesRef.current = candlestickSeries

    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
      chart.remove()
      seriesRef.current = null
    }
  }, [])

  // 2. Set/Update historical data when chartData changes
  useEffect(() => {
    if (seriesRef.current) {
      seriesRef.current.setData(chartData)
    }
  }, [chartData])

  // 3. Handle real-time updates from WebSocket
  useEffect(() => {
    if (!lastTrade || !seriesRef.current) return

    try {
      const ethVal = parseFloat(formatEther(BigInt(lastTrade.eth_amount)))
      const tokenVal = parseFloat(formatEther(BigInt(lastTrade.token_amount)))
      if (tokenVal === 0) return

      const price = ethVal / tokenVal

      let durationSeconds = 60
      if (resolution === "5m") durationSeconds = 300
      if (resolution === "15m") durationSeconds = 900
      if (resolution === "1h") durationSeconds = 3600
      if (resolution === "1d") durationSeconds = 86400

      const tradeTime = Math.floor(new Date(lastTrade.timestamp).getTime() / 1000)
      const bucketTime = Math.floor(tradeTime / durationSeconds) * durationSeconds

      seriesRef.current.update({
        time: bucketTime,
        open: price,
        high: price,
        low: price,
        close: price,
      })
    } catch (e) {
      console.error("Realtime chart update failed:", e)
    }
  }, [lastTrade, resolution])

  return (
    <div className="glass-panel" style={{ padding: '16px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h3 style={{ fontSize: '1rem', color: '#fff', margin: 0 }}>Courbe de Prix</h3>
        <div style={{ display: 'flex', gap: '4px', background: 'rgba(255,255,255,0.03)', padding: '2px', borderRadius: '8px', border: '1px solid var(--border-glass)' }}>
          {["1m", "5m", "1h", "1d"].map((res) => (
            <button
              key={res}
              onClick={() => setResolution(res)}
              style={{
                background: resolution === res ? 'var(--color-primary)' : 'transparent',
                border: 'none',
                color: '#fff',
                fontSize: '0.75rem',
                padding: '4px 8px',
                borderRadius: '6px',
                cursor: 'pointer',
                fontWeight: resolution === res ? 'bold' : 'normal',
                transition: 'all 0.2s ease',
              }}
            >
              {res}
            </button>
          ))}
        </div>
      </div>
      <div ref={chartContainerRef} style={{ width: '100%', position: 'relative' }} />
    </div>
  )
}
