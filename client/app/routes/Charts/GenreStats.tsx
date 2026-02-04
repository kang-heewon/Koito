import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  ResponsiveContainer,
  Treemap,
  Tooltip,
  type TooltipProps,
} from "recharts";
import { getGenreStats, type GenreStatsResponse } from "api/api";
import PeriodSelector from "~/components/PeriodSelector";
import * as Tabs from "@radix-ui/react-tabs";

const COLORS = [
  "#8884d8",
  "#83a6ed",
  "#8dd1e1",
  "#82ca9d",
  "#a4de6c",
  "#d0ed57",
  "#ffc658",
  "#ff8042",
  "#ffbb28",
  "#0088fe",
  "#00c49f",
  "#ffbb28",
  "#ff8042",
];

const getColor = (index: number) => COLORS[index % COLORS.length];

export default function GenreStats() {
  const [period, setPeriod] = useState("month");
  const [metric, setMetric] = useState<"count" | "time">("count");

  const { data, isLoading, error } = useQuery<GenreStatsResponse>({
    queryKey: ["genre-stats", period, metric],
    queryFn: () => getGenreStats(period, metric),
  });

  const chartData =
    data?.stats.map((stat, index) => ({
      ...stat,
      fill: getColor(index),
    })) || [];

  const treeData = [
    {
      name: "Genres",
      children: chartData,
    },
  ];

  return (
    <main className="flex flex-grow justify-center pb-4">
      <div className="flex-1 flex flex-col items-center gap-16 min-h-0 mt-20">
        <h1 className="text-2xl font-bold">장르 통계</h1>

        <div className="w-full max-w-7xl px-5 flex flex-col gap-10">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <PeriodSelector current={period} setter={setPeriod} />

            <Tabs.Root
              value={metric}
              onValueChange={(v) => setMetric(v as "count" | "time")}
              className="flex bg-(--color-bg-secondary) rounded-lg p-1"
            >
              <Tabs.List className="flex gap-1">
                <Tabs.Trigger
                  value="count"
                  className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${
                    metric === "count"
                      ? "bg-(--color-bg) shadow-sm text-(--color-fg)"
                      : "text-(--color-fg-secondary) hover:text-(--color-fg)"
                  }`}
                >
                  청취 횟수
                </Tabs.Trigger>
                <Tabs.Trigger
                  value="time"
                  className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${
                    metric === "time"
                      ? "bg-(--color-bg) shadow-sm text-(--color-fg)"
                      : "text-(--color-fg-secondary) hover:text-(--color-fg)"
                  }`}
                >
                  청취 시간
                </Tabs.Trigger>
              </Tabs.List>
            </Tabs.Root>
          </div>

          <div className="h-[600px] w-full rounded-xl overflow-hidden">
            {isLoading ? (
              <div className="w-full h-full flex items-center justify-center">
                Loading...
              </div>
            ) : error ? (
              <div className="w-full h-full flex items-center justify-center text-red-500">
                Error loading data
              </div>
            ) : chartData.length === 0 ? (
              <div className="w-full h-full flex items-center justify-center text-(--color-fg-secondary)">
                데이터가 없습니다
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <Treemap
                  data={treeData}
                  dataKey="value"
                  aspectRatio={4 / 3}
                  stroke="#fff"
                  fill="#8884d8"
                  content={<CustomContent />}
                >
                  <Tooltip content={<CustomTooltip metric={metric} />} />
                </Treemap>
              </ResponsiveContainer>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}

const CustomContent = (props: any) => {
  const { x, y, width, height, name, value, fill } = props;

  if (width < 50 || height < 30) return null;

  return (
    <g>
      <rect
        x={x}
        y={y}
        width={width}
        height={height}
        style={{
          fill: fill,
          stroke: "var(--color-bg)",
          strokeWidth: 2,
        }}
      />
      <text
        x={x + width / 2}
        y={y + height / 2}
        textAnchor="middle"
        dominantBaseline="middle"
        fill="#fff"
        fontSize={14}
        fontWeight="bold"
        style={{ textShadow: "0 1px 2px rgba(0,0,0,0.3)" }}
      >
        {name}
      </text>
    </g>
  );
};

const CustomTooltip = ({
  active,
  payload,
  metric,
}: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    const value = data.value;

    const formattedValue =
      metric === "time"
        ? `${Math.floor(value / 60)}시간 ${value % 60}분`
        : `${value.toLocaleString()}회`;

    return (
      <div className="bg-(--color-bg) border border-(--color-border) p-3 rounded shadow-lg">
        <p className="font-bold mb-1">{data.name}</p>
        <p className="text-sm text-(--color-fg-secondary)">
          {formattedValue}
        </p>
      </div>
    );
  }
  return null;
};
