/** @jsxImportSource react */
import { useCallback, useEffect, useState } from 'react';

export default function TerminalTyping() {
  const [text, setText] = useState('');
  const [showOutput, setShowOutput] = useState(false);
  const [outputText, setOutputText] = useState('');
  const [showNewPrompt, setShowNewPrompt] = useState(false);
  const [secondCommand, setSecondCommand] = useState('');
  const [key, setKey] = useState(0);
  
  const command = 'github.com/ChefBingbong/viem-go';
  const fullOutput = `go: downloading github.com/ChefBingbong/viem-go v1.0.0
go: upgraded go 1.21 => 1.24.0
go: added github.com/ChefBingbong/viem-go v1.0.0`;
  const secondCmd = 'go run main.go';

  const resetAnimation = useCallback(() => {
    setText('');
    setShowOutput(false);
    setOutputText('');
    setShowNewPrompt(false);
    setSecondCommand('');
    setKey(k => k + 1);
  }, []);

  useEffect(() => {
    let i = 0;
    let outputComplete = false;
    
    // Type the command
    const typeInterval = setInterval(() => {
      if (i < command.length) {
        setText(command.slice(0, i + 1));
        i++;
      } else {
        clearInterval(typeInterval);
        // Delay before showing output
        setTimeout(() => {
          setShowOutput(true);
          // Type output very fast (almost instant)
          let j = 0;
          const outputInterval = setInterval(() => {
            if (j < fullOutput.length) {
              // Type 10 characters at a time for rapid effect
              const chunkSize = 10;
              j = Math.min(j + chunkSize, fullOutput.length);
              setOutputText(fullOutput.slice(0, j));
            } else {
              clearInterval(outputInterval);
              if (!outputComplete) {
                outputComplete = true;
                // Show new prompt after a small delay
                setTimeout(() => {
                  setShowNewPrompt(true);
                  // Start typing second command after a short delay
                  setTimeout(() => {
                    let k = 0;
                    const secondTypeInterval = setInterval(() => {
                      if (k < secondCmd.length) {
                        setSecondCommand(secondCmd.slice(0, k + 1));
                        k++;
                      } else {
                        clearInterval(secondTypeInterval);
                        // Wait 4 seconds then restart
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
    <div className="terminal-typing-content">
      <div className="terminal-line">
        <span className="terminal-prompt">user@golem</span>
        <span className="terminal-dir"> ~</span>
        <span className="terminal-symbol"> % </span>
        <span className="terminal-cmd">go get </span>
        <span>{text}</span>
        {!showOutput && <span className="typing-cursor">|</span>}
      </div>
      {showOutput && (
        <div className="terminal-output-lines">
          {outputText.split('\n').map((line, idx) => (
            <div key={idx} className="output-line">{line}</div>
          ))}
        </div>
      )}
      {showNewPrompt && (
        <div className="terminal-line new-prompt">
          <span className="terminal-prompt">user@golem</span>
          <span className="terminal-dir"> ~</span>
          <span className="terminal-symbol"> % </span>
          <span>{secondCommand}</span>
          <span className="typing-cursor">|</span>
        </div>
      )}
    </div>
  );
}
