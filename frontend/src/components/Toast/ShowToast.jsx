import { Toast, ToastContainer } from "react-bootstrap";

export default function ShowToast({ toastMsgObj, setShowToast }) {
  return (
    <ToastContainer position="top-end" className="p-3">
      <Toast
        onClose={() => setShowToast(false)}
        show={true}
        bg={
          toastMsgObj.variant.toLowerCase() === "error" ? "danger" : "success"
        }
        delay={toastMsgObj.variant === "Error" ? 5000 : 3000}
        autohide={toastMsgObj.variant !== "Error"}
      >
        <Toast.Header>
          <strong className="me-auto">{toastMsgObj.variant + "!"}</strong>
        </Toast.Header>
        <Toast.Body style={{ color: "antiquewhite" }}>
          <pre>{toastMsgObj.message.trim()}</pre>
        </Toast.Body>
      </Toast>
    </ToastContainer>
  );
}
