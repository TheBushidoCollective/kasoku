'use client'

import { useEffect, useState } from 'react'

interface LiveCounterProps {
  initialValue: number
  incrementPerSecond: number
  suffix?: string
  prefix?: string
  decimals?: number
}

export function LiveCounter({
  initialValue,
  incrementPerSecond,
  suffix = '',
  prefix = '',
  decimals = 0
}: LiveCounterProps) {
  const [count, setCount] = useState(initialValue)
  const [hasAnimated, setHasAnimated] = useState(false)

  useEffect(() => {
    // Initial animation to the starting value
    if (!hasAnimated) {
      let startTime: number
      let animationFrame: number
      const duration = 2000 // 2 seconds for initial animation

      const animate = (timestamp: number) => {
        if (!startTime) startTime = timestamp
        const progress = timestamp - startTime
        const percentage = Math.min(progress / duration, 1)

        // Easing function for smooth animation
        const easeOutQuart = 1 - Math.pow(1 - percentage, 4)
        const current = easeOutQuart * initialValue

        setCount(current)

        if (percentage < 1) {
          animationFrame = requestAnimationFrame(animate)
        } else {
          setCount(initialValue)
          setHasAnimated(true)
        }
      }

      animationFrame = requestAnimationFrame(animate)

      return () => {
        if (animationFrame) {
          cancelAnimationFrame(animationFrame)
        }
      }
    }
  }, [initialValue, hasAnimated])

  useEffect(() => {
    // After initial animation, increment every second
    if (!hasAnimated) return

    const interval = setInterval(() => {
      setCount(prev => prev + incrementPerSecond)
    }, 1000)

    return () => clearInterval(interval)
  }, [incrementPerSecond, hasAnimated])

  const formattedCount = decimals > 0
    ? count.toFixed(decimals)
    : Math.floor(count).toLocaleString()

  return (
    <span className="tabular-nums">
      {prefix}{formattedCount}{suffix}
    </span>
  )
}
