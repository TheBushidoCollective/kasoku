'use client'

import { useEffect, useState } from 'react'

interface StatCounterProps {
  end: number
  duration?: number
  suffix?: string
  prefix?: string
  decimals?: number
}

export function StatCounter({
  end,
  duration = 2000,
  suffix = '',
  prefix = '',
  decimals = 0
}: StatCounterProps) {
  const [count, setCount] = useState(0)

  useEffect(() => {
    let startTime: number
    let animationFrame: number

    const animate = (timestamp: number) => {
      if (!startTime) startTime = timestamp
      const progress = timestamp - startTime
      const percentage = Math.min(progress / duration, 1)

      // Easing function for smooth animation
      const easeOutQuart = 1 - Math.pow(1 - percentage, 4)
      const current = easeOutQuart * end

      setCount(current)

      if (percentage < 1) {
        animationFrame = requestAnimationFrame(animate)
      } else {
        setCount(end)
      }
    }

    animationFrame = requestAnimationFrame(animate)

    return () => {
      if (animationFrame) {
        cancelAnimationFrame(animationFrame)
      }
    }
  }, [end, duration])

  const formattedCount = decimals > 0
    ? count.toFixed(decimals)
    : Math.floor(count).toLocaleString()

  return (
    <span className="tabular-nums">
      {prefix}{formattedCount}{suffix}
    </span>
  )
}
