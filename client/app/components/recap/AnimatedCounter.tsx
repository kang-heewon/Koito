import { useEffect, useRef, useState } from "react";
import {
  animate,
  motion,
  useInView,
  useMotionValue,
  useMotionValueEvent,
  useTransform,
} from "motion/react";

interface AnimatedCounterProps {
  value: number;
  duration?: number;
  suffix?: string;
  prefix?: string;
}

function formatCounterValue(value: number, prefix: string, suffix: string) {
  const maximumFractionDigits = Number.isInteger(value) ? 0 : 1;

  return `${prefix}${value.toLocaleString(undefined, { maximumFractionDigits })}${suffix}`;
}

export default function AnimatedCounter({
  value,
  duration = 1.2,
  suffix = "",
  prefix = "",
}: AnimatedCounterProps) {
  const ref = useRef<HTMLSpanElement | null>(null);
  const isInView = useInView(ref, { once: true, amount: 0.7 });
  const motionValue = useMotionValue(0);
  const transformedValue = useTransform(motionValue, (latest) =>
    formatCounterValue(latest, prefix, suffix),
  );
  const [displayValue, setDisplayValue] = useState(formatCounterValue(0, prefix, suffix));

  useMotionValueEvent(transformedValue, "change", (latest) => {
    setDisplayValue(latest);
  });

  useEffect(() => {
    if (!isInView) {
      return;
    }

    motionValue.set(0);

    const controls = animate(motionValue, value, {
      duration,
      ease: "easeOut",
    });

    return () => {
      controls.stop();
    };
  }, [duration, isInView, motionValue, value]);

  return (
    <motion.span ref={ref} className="tabular-nums">
      {displayValue}
    </motion.span>
  );
}
