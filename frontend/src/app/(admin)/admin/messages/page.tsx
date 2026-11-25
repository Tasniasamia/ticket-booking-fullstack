import type { Metadata } from "next";
import { MessagesPanel } from "@/components/messaging/messages-panel";
import { PageHeader } from "@/components/dashboard/shared/page-header";

export const metadata: Metadata = { title: "Messages" };

export default function AdminMessagesPage() {
  return (
    <div>
      <PageHeader
        title="Messages"
        subtitle="Message any user. Start a chat by entering their user ID, or open an existing conversation."
      />
      <MessagesPanel allowStartById className="h-[calc(100dvh-14.5rem)] min-h-[26rem]" />
    </div>
  );
}
