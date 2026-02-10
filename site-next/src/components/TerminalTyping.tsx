"use client";

import { useCallback, useEffect, useState } from "react";

export default function TerminalTyping() {
  const [text, setText] = useState("");
  const [showOutput, setShowOutput] = useState(false);
  const [outputText, setOutputText] = useState("");
  const [showNewPrompt, setShowNewPrompt] = useState(false);
  const [secondCommand, setSecondCommand] = useState("");
  const [key, setKey] = useState(0);

  const command = "github.com/ChefBingbong/viem-go";
  const fullOutput = `go: downloading github.com/ChefBingbong/viem-go v1.0.0
go: upgraded go 1.21 => 1.24.0
go: added github.com/ChefBingbong/viem-go v1.0.0`;
  const secondCmd = "go run main.go";

  const resetAnimation = useCallback(() => {
    setText("");
    setShowOutput(false);
    setOutputText("");
    setShowNewPrompt(false);
    setSecondCommand("");
    setKey((k) => k + 1);
  }, []);

  useEffect(() => {
    let i = 0;
    let outputComplete = false;

    const typeInterval = setInterval(() => {
      if (i < command.length) {
        setText(command.slice(0, i + 1));
        i++;
      } else {
        clearInterval(typeInterval);
        setTimeout(() => {
          setShowOutput(true);
          let j = 0;
          const outputInterval = setInterval(() => {
            if (j < fullOutput.length) {
              const chunkSize = 10;
              j = Math.min(j + chunkSize, fullOutput.length);
              setOutputText(fullOutput.slice(0, j));
            } else {
              clearInterval(outputInterval);
              if (!outputComplete) {
                outputComplete = true;
                setTimeout(() => {
                  setShowNewPrompt(true);
                  setTimeout(() => {
                    let k = 0;
                    const secondTypeInterval = setInterval(() => {
                      if (k < secondCmd.length) {
                        setSecondCommand(secondCmd.slice(0, k + 1));
                        k++;
                      } else {
                        clearInterval(secondTypeInterval);
                        setTimeout(() => {
                          resetAnimation();
                        }, 4000);
                      }
                    }, 50);
                  }, 500);
                }, 300);
              }
            }
          }, 5);
        }, 800);
      }
    }, 40);

    return () => clearInterval(typeInterval);
  }, [key, resetAnimation]);

  return (
    <div className="font-mono text-[0.9375rem] text-foreground-secondary w-full h-56">
      <div className="whitespace-nowrap">
        <span className="text-terminal-user">user@golem</span>
        <span className="text-terminal-path"> ~</span>
        <span className="text-foreground-muted"> % </span>
        <span className="text-terminal-command">go get </span>
        <span>{text}</span>
        {!showOutput && (
          <span className="inline-block text-terminal-cursor animate-cursor-blink ml-px">
            |
          </span>
        )}
      </div>
      {showOutput && (
        <div className="mt-1.5 text-[0.8rem] text-terminal-output leading-normal">
          {outputText.split("\n").map((line, idx) => (
            <div key={idx} className="py-px">
              {line}
            </div>
          ))}
        </div>
      )}
      {showNewPrompt && (
        <div className="whitespace-nowrap mt-3">
          <span className="text-terminal-user">user@golem</span>
          <span className="text-terminal-path"> ~</span>
          <span className="text-foreground-muted"> % </span>
          <span>{secondCommand}</span>
          <span className="inline-block text-terminal-cursor animate-cursor-blink ml-px">
            |
          </span>
        </div>
      )}
    </div>
  );
}
