import "./ErrorMessage.css";

function ErrorMessage({ message, onClose }) {
  if (!message) return null;
  return (
    <div className="error-banner" role="alert">
      <span className="error-icon">⚠️</span>
      <span className="error-text">{message}</span>
      {onClose && (
        <button
          onClick={onClose}
          className="error-close-btn"
          aria-label="Fechar erro"
        >
          &times;
        </button>
      )}
    </div>
  );
}

export default ErrorMessage;
