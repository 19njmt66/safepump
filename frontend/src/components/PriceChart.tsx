import { useEffect, useRef, useState } from 'react'
import { createChart, ColorType, CandlestickSeries, HistogramSeries } from 'lightweight-charts'
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
  volume?: number
}

export default function PriceChart({ tokenAddress, lastTrade }: PriceChartProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const seriesRef = useRef<any>(null)
  const volumeSeriesRef = useRef<any>(null)
  const [resolution, setResolution] = useState("1m")
  const [chartData, setChartData] = useState<CandleData[]>([])
  
  // Real-time hover state
  const [hoverData, setHoverData] = useState<{ open: number, high: number, low: number, close: number, volume: number } | null>(null)

  // Fetch historical candles
  const fetchCandles = async () => {
    try {
      const hostname = typeof window !== 'undefined' ? window.location.hostname : 'localhost';
      const res = await fetch(`http://${hostname}:8080/api/v1/tokens/${tokenAddress}/candles?resolution=${resolution}`)
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

  // Initialize the chart on mount
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
      height: 450,
      timeScale: {
        timeVisible: true,
        secondsVisible: false,
      },
    })

    // Add Candlestick Series
    const candlestickSeries = chart.addSeries(CandlestickSeries, {
      upColor: '#10b981',
      downColor: '#ef4444',
      borderVisible: false,
      wickUpColor: '#10b981',
      wickDownColor: '#ef4444',
      priceFormat: {
        type: 'price',
        precision: 8,
        minMove: 0.00000001,
      },
    })

    // Add Volume Series overlay
    const volumeSeries = chart.addSeries(HistogramSeries, {
      color: 'rgba(59, 130, 246, 0.2)',
      priceFormat: {
        type: 'volume',
      },
      priceScaleId: '', // overlay mode
    })

    volumeSeries.priceScale().applyOptions({
      scaleMargins: {
        top: 0.8, // volume bars occupy only the bottom 20% of chart
        bottom: 0,
      },
    })

    seriesRef.current = candlestickSeries
    volumeSeriesRef.current = volumeSeries

    // Subscribe to crosshair moves to show real-time legend
    chart.subscribeCrosshairMove((param) => {
      if (
        param.point === undefined ||
        !param.time ||
        param.point.x < 0 ||
        param.point.y < 0
      ) {
        setHoverData(null)
      } else {
        const candle = param.seriesData.get(candlestickSeries) as any
        const volume = param.seriesData.get(volumeSeries) as any
        if (candle) {
          setHoverData({
            open: candle.open,
            high: candle.high,
            low: candle.low,
            close: candle.close,
            volume: volume ? volume.value : 0
          })
        }
      }
    })

    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
      chart.remove()
      seriesRef.current = null
      volumeSeriesRef.current = null
    }
  }, [])

  useEffect(() => {
    if (seriesRef.current) {
      seriesRef.current.setData(chartData)
    }
    if (volumeSeriesRef.current) {
      const volumeData = chartData.map(c => ({
        time: c.time,
        value: c.volume || 0,
        color: c.close >= c.open ? 'rgba(16, 185, 129, 0.2)' : 'rgba(239, 68, 68, 0.2)'
      }))
      volumeSeriesRef.current.setData(volumeData)
    }
  }, [chartData])

  // Handle real-time updates from WebSocket
  useEffect(() => {
    if (!lastTrade || !seriesRef.current) return

    try {
      const ethVal = parseFloat(formatEther(BigInt(lastTrade.eth_amount)))
      const tokenVal = parseFloat(formatEther(BigInt(lastTrade.token_amount)))
      if (tokenVal === 0) return

      const price = ethVal / tokenVal
      const volumeTokens = tokenVal / 1e18

      let durationSeconds = 60
      if (resolution === "5m") durationSeconds = 300
      if (resolution === "15m") durationSeconds = 900
      if (resolution === "1h") durationSeconds = 3600
      if (resolution === "1d") durationSeconds = 86400

      const tradeTime = Math.floor(new Date(lastTrade.timestamp).getTime() / 1000)
      const bucketTime = Math.floor(tradeTime / durationSeconds) * durationSeconds

      seriesRef.current.update({
        time: bucketTime as any,
        open: price,
        high: price,
        low: price,
        close: price,
      })

      if (volumeSeriesRef.current) {
        volumeSeriesRef.current.update({
          time: bucketTime as any,
          value: volumeTokens,
          color: lastTrade.is_buy ? 'rgba(16, 185, 129, 0.2)' : 'rgba(239, 68, 68, 0.2)'
        })
      }
    } catch (e) {
      console.error("Realtime chart update failed:", e)
    }
  }, [lastTrade, resolution])

  return (
    <div className="glass-panel" style={{ padding: '16px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flexWrap: 'wrap' }}>
          <h3 style={{ fontSize: '1rem', color: '#fff', margin: 0 }}>Courbe de Prix</h3>
          {hoverData && (
            <div style={{ display: 'flex', gap: '8px', fontSize: '0.75rem', fontFamily: 'monospace', color: 'var(--color-text-muted)' }}>
              <span>O: <span style={{ color: hoverData.close >= hoverData.open ? 'var(--color-success)' : 'var(--color-danger)' }}>{hoverData.open.toFixed(8)}</span></span>
              <span>H: <span style={{ color: hoverData.close >= hoverData.open ? 'var(--color-success)' : 'var(--color-danger)' }}>{hoverData.high.toFixed(8)}</span></span>
              <span>L: <span style={{ color: hoverData.close >= hoverData.open ? 'var(--color-success)' : 'var(--color-danger)' }}>{hoverData.low.toFixed(8)}</span></span>
              <span>C: <span style={{ color: hoverData.close >= hoverData.open ? 'var(--color-success)' : 'var(--color-danger)' }}>{hoverData.close.toFixed(8)}</span></span>
              <span>V: <span style={{ color: 'var(--color-primary)' }}>{hoverData.volume.toLocaleString(undefined, { maximumFractionDigits: 0 })}</span></span>
            </div>
          )}
        </div>
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
