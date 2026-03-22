import { type ReactNode, useRef } from "react";
import { motion, useInView } from "motion/react";
import AnimatedCounter from "./AnimatedCounter";

interface StatRevealProps {
  value: string | number;
  label: string;
  icon?: ReactNode;
}

const easing = [0.22, 1, 0.36, 1] as const;

export default function StatReveal({ value, label, icon }: StatRevealProps) {
  const ref = useRef<HTMLDivElement | null>(null);
  const isInView = useInView(ref, { once: true, amount: 0.4 });

  return (
    <motion.div
      ref={ref}
      className="flex min-h-64 flex-col justify-between gap-8 rounded-[28px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/80 px-6 py-7 text-[var(--color-fg)] backdrop-blur-sm sm:px-8 sm:py-9"
      initial={false}
      animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 32 }}
      transition={{ duration: 0.7, ease: easing }}
    >
      <div className="flex items-center gap-3 text-sm font-semibold uppercase tracking-[0.24em] text-[var(--color-primary)]">
        {icon ? <span className="flex h-10 w-10 items-center justify-center rounded-full border border-[var(--color-primary)]/20 bg-[var(--color-bg)]/70 text-lg">{icon}</span> : null}
        <span>{label}</span>
      </div>

      <div className="header-font text-5xl font-semibold leading-none tracking-[-0.04em] text-[var(--color-fg)] sm:text-6xl">
        {typeof value === "number" ? <AnimatedCounter value={value} /> : value}
      </div>
    </motion.div>
  );
}
