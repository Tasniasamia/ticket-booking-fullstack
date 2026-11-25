import { Bone, WRAP } from "@/components/skeleton/sections";

export default function PublicLoading() {
  return (
    <div role="status" aria-busy="true" aria-label="Loading" className="space-y-12 py-10">
      {/* Hero Skeleton */}
      <section className={`${WRAP} space-y-4 text-center`}>
        <Bone className="mx-auto h-12 w-3/4 max-w-lg" />
        <Bone className="mx-auto h-4 w-full max-w-md" />
        <Bone className="mx-auto h-4 w-2/3 max-w-md" />
      </section>

      {/* Content Skeleton */}
      <section className={`${WRAP} grid gap-6 sm:grid-cols-2 lg:grid-cols-3`}>
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="space-y-3 rounded-3xl border border-fg/10 bg-fg/5 p-3">
            <Bone className="aspect-[16/10] w-full rounded-2xl" />
            <div className="space-y-2 p-2">
              <Bone className="h-4 w-20" />
              <Bone className="h-6 w-4/5" />
              <Bone className="h-4 w-full" />
            </div>
          </div>
        ))}
      </section>
    </div>
  );
}
