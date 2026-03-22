import { type ReactNode, useRef } from "react";
import { motion, useInView } from "motion/react";

interface ScrollSectionProps {
  children: ReactNode;
  className?: string;
  delay?: number;
}

const easing = [0.22, 1, 0.36, 1] as const;

export default function ScrollSection({ children, className, delay = 0 }: ScrollSectionProps) {
  const ref = useRef<HTMLElement | null>(null);
  const isInView = useInView(ref, { once: true, amount: 0.25 });
  const sectionClassName = ["min-h-screen flex items-center justify-center", className]
    .filter(Boolean)
    .join(" ");

  return (
    <section ref={ref} className={sectionClassName}>
      <motion.div
        className="w-full"
        initial={false}
        animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 48 }}
        transition={{ duration: 0.8, delay, ease: easing }}
      >
        {children}
      </motion.div>
    </section>
  );
}
