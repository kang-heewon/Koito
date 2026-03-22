import { useRef } from "react";
import { motion, useInView } from "motion/react";

interface TopItemCardProps {
  rank: number;
  name: string;
  imageUrl?: string;
  subtitle?: string;
  plays?: number;
}

const easing = [0.22, 1, 0.36, 1] as const;

export default function TopItemCard({ rank, name, imageUrl, subtitle, plays }: TopItemCardProps) {
  const ref = useRef<HTMLDivElement | null>(null);
  const isInView = useInView(ref, { once: true, amount: 0.35 });
  const delay = Math.max(rank - 1, 0) * 0.12;

  return (
    <motion.div
      ref={ref}
      className="relative flex items-center gap-4 rounded-[28px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/80 p-4 text-[var(--color-fg)] backdrop-blur-sm sm:gap-5 sm:p-5"
      initial={false}
      animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 28 }}
      transition={{ duration: 0.6, delay, ease: easing }}
    >
      <div className="absolute left-4 top-4 flex h-9 min-w-9 items-center justify-center rounded-full border border-[var(--color-primary)]/20 bg-[var(--color-bg)]/85 px-2 text-sm font-semibold text-[var(--color-primary)] sm:left-5 sm:top-5">
        #{rank}
      </div>

      <div className="mt-10 flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-[22px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)] text-3xl font-semibold text-[var(--color-primary)] sm:mt-0 sm:h-24 sm:w-24">
        {imageUrl ? (
          <img src={imageUrl} alt={name} className="h-full w-full object-cover" />
        ) : (
          <span className="header-font leading-none">{name.charAt(0).toUpperCase()}</span>
        )}
      </div>

      <div className="min-w-0 flex-1 pt-10 sm:pt-0">
        <div className="header-font truncate text-2xl font-semibold tracking-[-0.03em] text-[var(--color-fg)] sm:text-3xl">
          {name}
        </div>

        {subtitle ? (
          <div className="mt-1 truncate text-sm text-[var(--color-fg)]/72 sm:text-base">
            {subtitle}
          </div>
        ) : null}

        {typeof plays === "number" ? (
          <div className="mt-4 inline-flex items-center rounded-full border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/75 px-3 py-1 text-sm font-medium text-[var(--color-primary)] sm:text-base">
            {plays.toLocaleString()} plays
          </div>
        ) : null}
      </div>
    </motion.div>
  );
}
