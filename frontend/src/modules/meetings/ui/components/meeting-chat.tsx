import { useState, useRef, useEffect } from "react";
import { SendIcon, SparklesIcon, Loader2Icon, User2Icon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import { authClient } from "@/lib/auth-client";
import type { Meeting } from "../../types";
import { GeneratedAvatar } from "@/components/generated-avatar";
import ReactMarkdown from "react-markdown";

interface MeetingChatProps {
  meeting: Meeting;
}

interface Message {
  id: string;
  role: "user" | "ai";
  content: string;
  timestamp: Date;
}

// Compact markdown components for chat
const chatMarkdownComponents = {
  p: ({ node, ...props }: any) => (
    <p className="mb-2 last:mb-0 leading-relaxed" {...props} />
  ),
  ul: ({ node, ...props }: any) => (
    <ul className="list-disc list-outside ml-4 space-y-1 mb-2" {...props} />
  ),
  ol: ({ node, ...props }: any) => (
    <ol className="list-decimal list-outside ml-4 space-y-1 mb-2" {...props} />
  ),
  li: ({ node, ...props }: any) => <li className="pl-1" {...props} />,
  strong: ({ node, ...props }: any) => (
    <span className="font-semibold" {...props} />
  ),
  code: ({ node, ...props }: any) => (
    <code
      className="bg-black/10 dark:bg-white/10 px-1 py-0.5 rounded text-xs font-mono"
      {...props}
    />
  ),
};

export const MeetingChat = ({ meeting }: MeetingChatProps) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages, isStreaming]);

  // Fetch chat history on mount
  useEffect(() => {
    const fetchHistory = async () => {
      try {
        const token = await authClient.token();
        const serverUrl = import.meta.env.VITE_SERVER_URL;

        const response = await fetch(`${serverUrl}/chat/${meeting.id}`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token.data?.token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch chat history");
        }

        const data = await response.json();

        const historyMessages: Message[] = data.messages.map((msg: any) => ({
          id: msg.id,
          role: msg.role as "user" | "ai",
          content: msg.content,
          timestamp: new Date(msg.createdAt.Time || msg.createdAt),
        }));

        setMessages(historyMessages);
      } catch (error) {
        console.error("Error fetching chat history:", error);
      }
    };

    fetchHistory();
  }, [meeting.id]);

  const handleSendMessage = async () => {
    if (!inputValue.trim() || isStreaming) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: "user",
      content: inputValue,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputValue("");
    setIsStreaming(true);

    const aiMessageId = (Date.now() + 1).toString();
    setMessages((prev) => [
      ...prev,
      {
        id: aiMessageId,
        role: "ai",
        content: "",
        timestamp: new Date(),
      },
    ]);

    try {
      const token = await authClient.token();
      const serverUrl = import.meta.env.VITE_SERVER_URL;

      const response = await fetch(`${serverUrl}/chat/${meeting.id}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token.data?.token}`,
        },
        body: JSON.stringify({ message: inputValue }),
      });

      if (!response.ok) {
        throw new Error("Failed to send message");
      }

      if (!response.body) {
        throw new Error("No response body");
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        buffer += chunk;

        // Process complete lines (split by newline)
        const lines = buffer.split("\n");
        // Keep the last incomplete line in the buffer
        buffer = lines.pop() || "";

        for (const line of lines) {
          // Check if this line starts with "data:"
          if (line.startsWith("data:")) {
            // Extract everything after "data:" - preserve all spacing as-is
            const content = line.substring(5);

            // Append the content exactly as received from server
            setMessages((prev) =>
              prev.map((msg) =>
                msg.id === aiMessageId
                  ? { ...msg, content: msg.content + content }
                  : msg
              )
            );
          }
        }
      }
    } catch (error) {
      console.error("Error sending message:", error);
      // Optionally handle error state in UI
      setMessages((prev) =>
        prev.map((msg) =>
          msg.id === aiMessageId
            ? {
                ...msg,
                content: msg.content + "\n[Error: Failed to get response]",
              }
            : msg
        )
      );
    } finally {
      setIsStreaming(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="flex flex-col h-full w-[70%] mx-auto">
      <div className="flex flex-col overflow-y-auto gap-6 py-6 px-4">
        {messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-[400px] text-center space-y-6 opacity-60">
            <div className="bg-primary/5 p-5 rounded-full ring-1 ring-primary/10">
              <SparklesIcon className="w-8 h-8 text-primary" />
            </div>
            <div className="space-y-2 max-w-[280px]">
              <h3 className="font-medium text-lg text-foreground">
                Ask AI about this meeting
              </h3>
              <p className="text-sm text-muted-foreground leading-relaxed">
                Get summaries, action items, or clarify details directly from
                the transcript.
              </p>
            </div>
          </div>
        ) : (
          messages.map((message) => (
            <div
              key={message.id}
              className={cn(
                "flex w-full gap-4 items-start animate-in fade-in slide-in-from-bottom-2 duration-300",
                message.role === "user" ? "justify-end" : "justify-start"
              )}
            >
              {message.role === "ai" && (
                <GeneratedAvatar
                  seed={meeting.agentDetails.name}
                  variant="botttsNeutral"
                  className="size-8 shrink-0 ring-1 ring-border/50"
                />
              )}

              <div
                className={cn(
                  "relative px-5 py-3.5 text-sm max-w-[60%] shadow-sm",
                  message.role === "user"
                    ? "bg-primary text-primary-foreground rounded-2xl rounded-tr-sm"
                    : "bg-muted text-card-foreground rounded-2xl rounded-tl-sm border border-border/40 "
                )}
              >
                {message.role === "ai" ? (
                  <div className=" prose prose-sm dark:prose-invert max-w-[85%] prose-p:leading-relaxed prose-pre:p-0 prose-pre:bg-transparent">
                    <ReactMarkdown components={chatMarkdownComponents}>
                      {message.content}
                    </ReactMarkdown>
                    {isStreaming && message.content === "" && (
                      <span className="flex items-center gap-1 text-muted-foreground animate-pulse">
                        <SparklesIcon className="w-3 h-3" /> Thinking...
                      </span>
                    )}
                  </div>
                ) : (
                  <p className="leading-relaxed">{message.content}</p>
                )}
              </div>

              {message.role === "user" && (
                <div className="size-8 rounded-full bg-muted/50 flex items-center justify-center border border-border/40 shrink-0">
                  <User2Icon className="size-4 text-muted-foreground" />
                </div>
              )}
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      <div className="pt-4 mt-auto">
        <div className="relative flex items-center group">
          <Input
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Ask a question about the meeting..."
            className="pr-12 py-6 pl-6 rounded-full shadow-sm border-border/40 bg-muted/20 focus-visible:bg-background focus-visible:ring-1 focus-visible:ring-primary/20 focus-visible:border-primary/30 transition-all"
            disabled={isStreaming}
          />
          <Button
            size="icon"
            variant="ghost"
            className={cn(
              "absolute right-2 h-9 w-9 rounded-full transition-all duration-300",
              inputValue.trim()
                ? "bg-primary text-primary-foreground hover:bg-primary/90 shadow-md scale-100"
                : "text-muted-foreground hover:bg-muted scale-90 opacity-70"
            )}
            onClick={handleSendMessage}
            disabled={!inputValue.trim() || isStreaming}
          >
            {isStreaming ? (
              <Loader2Icon className="w-4 h-4 animate-spin" />
            ) : (
              <SendIcon className="w-4 h-4 ml-0.5" />
            )}
          </Button>
        </div>
        <div className="text-[10px] text-center text-muted-foreground/60 mt-3 font-medium tracking-wide uppercase">
          AI generated content may be inaccurate
        </div>
      </div>
    </div>
  );
};
