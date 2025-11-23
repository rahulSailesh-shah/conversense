export const markdownComponents = {
  h1: ({ node, ...props }: any) => (
    <h1
      className="text-3xl font-bold tracking-tight text-foreground mt-10 mb-6 first:mt-0"
      {...props}
    />
  ),
  h2: ({ node, ...props }: any) => (
    <h2
      className="text-2xl font-semibold tracking-tight text-foreground mt-8 mb-4 first:mt-0"
      {...props}
    />
  ),
  h3: ({ node, ...props }: any) => (
    <h3
      className="text-xl font-semibold tracking-tight text-foreground mt-8 mb-4 pb-2 border-b border-border/40 first:mt-0"
      {...props}
    />
  ),
  h4: ({ node, ...props }: any) => (
    <h4
      className="text-base font-semibold tracking-tight text-foreground/90 mt-6 mb-3"
      {...props}
    />
  ),
  ul: ({ node, ...props }: any) => (
    <ul
      className="list-disc list-outside ml-5 space-y-2 mb-6 marker:text-muted-foreground/60"
      {...props}
    />
  ),
  ol: ({ node, ...props }: any) => (
    <ol
      className="list-decimal list-outside ml-5 space-y-2 mb-6 marker:text-muted-foreground/60"
      {...props}
    />
  ),
  li: ({ node, ...props }: any) => (
    <li className="pl-1 leading-7 text-muted-foreground" {...props} />
  ),
  p: ({ node, ...props }: any) => (
    <p className="mb-4 last:mb-0 leading-7 text-muted-foreground" {...props} />
  ),
  strong: ({ node, ...props }: any) => (
    <span className="font-semibold text-foreground" {...props} />
  ),
  em: ({ node, ...props }: any) => (
    <span className="italic text-foreground/80" {...props} />
  ),
  a: ({ node, ...props }: any) => (
    <a
      className="text-primary font-medium hover:underline underline-offset-4 transition-colors"
      {...props}
    />
  ),
  blockquote: ({ node, ...props }: any) => (
    <blockquote
      className="border-l-4 border-primary/20 bg-muted/30 pl-4 py-2 pr-4 rounded-r italic my-6 text-muted-foreground"
      {...props}
    />
  ),
  code: ({ node, ...props }: any) => (
    <code
      className="bg-muted/50 px-1.5 py-0.5 rounded text-sm font-mono text-foreground/80 border border-border/50"
      {...props}
    />
  ),
};
