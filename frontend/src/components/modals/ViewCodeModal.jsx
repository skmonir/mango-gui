import { PrismLight as SyntaxHighlighter } from "react-syntax-highlighter";
import jsx from "react-syntax-highlighter/dist/esm/languages/prism/jsx";
import { darcula } from "react-syntax-highlighter/dist/esm/styles/prism";
import { useEffect, useState } from "react";
import { Button, Modal } from "react-bootstrap";
import DataService from "../../services/DataService.js";

SyntaxHighlighter.registerLanguage("jsx", jsx);

export default function ViewCodeModal({
  codePath,
  metadata,
  setShowCodeModal,
}) {
  const [code, setCode] = useState("");
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    if (codePath) {
      fetchCodeByPath(codePath);
    } else if (metadata) {
      fetchCodeByMetadata(metadata);
    }
  }, []);

  const fetchCodeByPath = (filepath) => {
    DataService.getCodeByPath({ filePath: filepath })
      .then((code) => {
        setCode(code);
        setShowModal(true);
      })
      .finally(() => setShowModal(true));
  };

  const fetchCodeByMetadata = (metadata) => {
    DataService.getCodeByMetadata(metadata)
      .then((code) => {
        setCode(code);
        setShowModal(true);
      })
      .finally(() => setShowModal(true));
  };

  const closeModal = () => {
    setShowModal(false);
    setTimeout(() => setShowCodeModal(false), 500);
  };

  return (
    <Modal
      show={showModal}
      onHide={closeModal}
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
      fullscreen={true}
    >
      <Modal.Header />
      <Modal.Body style={{ height: "80vh", overflowY: "auto" }}>
        <SyntaxHighlighter language="java" style={darcula}>
          {code}
        </SyntaxHighlighter>
      </Modal.Body>
      <Modal.Footer>
        <Button size="sm" variant="outline-danger" onClick={() => closeModal()}>
          Close Code
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
