import reactLogo from "../assets/react.svg";
import viteLogo from "../assets/vite.svg";
import golangLogo from "../assets/golang.svg";
import swcLogo from "../assets/swc.svg";
import "./TechStack.css";
import "../App.css";

function TechStack() {
  return (
    <div className="apple-subtitle tech-row">
      <div className="tech-item" title="React">
        <img
          src={reactLogo}
          className="tech-icon icon-react"
          alt="React Logo"
        />
        <span>React</span>
      </div>

      <span>+</span>

      <div className="tech-item" title="Vite">
        <img src={viteLogo} className="tech-icon icon-vite" alt="Vite Logo" />
        <span>Vite</span>
      </div>

      <span>+</span>

      <div className="tech-item" title="JavaScript">
        <img
          src={golangLogo}
          className="tech-icon icon-golang"
          alt="Golang Logo"
        />
        <span>Golang</span>
      </div>

      <span>+</span>

      <div className="tech-item" title="SWC">
        <img src={swcLogo} className="tech-icon icon-swc" alt="SWC Logo" />
        <span>SWC</span>
      </div>
    </div>
  );
}

export default TechStack;
