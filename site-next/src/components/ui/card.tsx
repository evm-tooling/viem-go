import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/cn";

const cardVariants = cva("rounded-xl border", {
  variants: {
    variant: {
      default: "border-accent/20 bg-gray-6/50",
      interactive:
        "border-accent/20 bg-gray-6/50 transition-all duration-200 hover:border-accent/40 hover:bg-gray-6/70 hover:-translate-y-0.5",
      surface: "border-card-border bg-card",
      surfaceInteractive:
        "border-card-border bg-card transition-all duration-200 hover:border-primary/30 hover:bg-card/80",
    },
    padding: {
      none: "",
      md: "p-6",
    },
  },
  defaultVariants: {
    variant: "default",
    padding: "md",
  },
});

export interface CardProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof cardVariants> {
  asChild?: boolean;
}

export function Card({
  className,
  variant,
  padding,
  asChild = false,
  ...props
}: CardProps) {
  const Comp = asChild ? Slot : "div";
  return (
    <Comp
      className={cn(cardVariants({ variant, padding }), className)}
      {...props}
    />
  );
}

export function CardTitle({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h3
      className={cn("heading-4 mb-1.5", className)}
      {...props}
    />
  );
}

export function CardDescription({
  className,
  ...props
}: React.HTMLAttributes<HTMLParagraphElement>) {
  return (
    <p className={cn("text-body", className)} {...props} />
  );
}

