import { type CSSProperties, type ReactNode } from "react";
import { motion } from "motion/react";

interface GradientBackgroundProps {
  colors: string[];
  children: ReactNode;
}

export default function GradientBackground({ colors, children }: GradientBackgroundProps) {
  const [first, second, third, fourth] = [
    colors[0] ?? "var(--color-primary)",
    colors[1] ?? "var(--color-bg)",
    colors[2] ?? colors[0] ?? "var(--color-primary)",
    colors[3] ?? "var(--color-bg)",
  ];

  const backgroundStyle: CSSProperties = {
    backgroundColor: "var(--color-bg)",
    backgroundImage: [
      `radial-gradient(circle at 18% 18%, ${first} 0%, transparent 42%)`,
      `radial-gradient(circle at 82% 24%, ${second} 0%, transparent 38%)`,
      `radial-gradient(circle at 50% 84%, ${third} 0%, transparent 44%)`,
      `linear-gradient(135deg, ${fourth} 0%, var(--color-bg) 100%)`,
    ].join(", "),
    backgroundPosition: "0% 50%",
    backgroundRepeat: "no-repeat",
    backgroundSize: "180% 180%",
  };

  return (
    <div className="relative overflow-hidden rounded-[32px] border border-[var(--color-primary)]/10 bg-[var(--color-bg)] text-[var(--color-fg)]">
      <motion.div
        aria-hidden
        className="absolute inset-0"
        style={backgroundStyle}
        animate={{
          backgroundPosition: ["0% 50%", "100% 50%", "0% 50%"],
          scale: [1, 1.04, 1],
        }}
        transition={{ duration: 18, ease: "linear", repeat: Number.POSITIVE_INFINITY }}
      />
      <div className="relative z-10">{children}</div>
    </div>
  );
}
