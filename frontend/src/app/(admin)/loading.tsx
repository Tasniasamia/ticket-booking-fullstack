import { Bone } from "@/components/skeleton/sections";

export default function AdminLoading() {
  return (
    <div className="space-y-6">
      {/* Header Skeleton */}
      <div className="space-y-2">
        <Bone className="h-8 w-48" />
        <Bone className="h-4 w-72" />
      </div>

      {/* Summary Cards Skeleton */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="rounded-2xl border border-fg/10 bg-panel/60 p-4 space-y-2">
            <Bone className="h-4 w-20" />
            <Bone className="h-7 w-28" />
          </div>
        ))}
      </div>

      {/* Controls & Table Skeleton */}
      <div className="space-y-4">
        <div className="flex gap-3">
          <Bone className="h-10 w-64 rounded-xl" />
          <Bone className="h-10 w-32 rounded-xl" />
        </div>
        <div className="rounded-2xl border border-fg/10 bg-panel/60 p-4 space-y-3">
          <Bone className="h-10 w-full rounded-lg" />
          {Array.from({ length: 5 }).map((_, i) => (
            <Bone key={i} className="h-12 w-full rounded-lg" />
          ))}
        </div>
      </div>
    </div>
  );
}
