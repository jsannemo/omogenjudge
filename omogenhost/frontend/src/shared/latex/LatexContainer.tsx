import React from 'react';
import renderMathInElement from "katex/dist/contrib/auto-render";
import 'katex/dist/katex.min.css';

export function LatexContainer({children}: { children: JSX.Element }) {
  const ref = React.useRef(null);
  React.useEffect(() => {
    if (ref.current) {
      ref.current.querySelectorAll('.tex2jax_process').forEach(
          (el: HTMLElement) => renderMathInElement(el, {
            delimiters: [
              {left: '$', right: '$', display: false},
              {left: '\\(', right: '\\)', display: false},
            ],
            ignoredClasses: ['*']
          }));
    }
  }, [ref.current]);
  return (
      <div ref={ref}>
        {children}
      </div>
  );
}