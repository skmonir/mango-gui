import { useEffect, useState } from "react";
import { Button, Modal } from "react-bootstrap";
import DataService from "../../services/DataService.js";
import Prism from "prismjs";
import "prismjs/themes/prism.css";
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import Editor from "react-simple-code-editor";

export default function ViewCodeModal({
  codePath,
  metadata,
  setShowCodeModal
}) {
  const [isCodeUpdated, setIsCodeUpdated] = useState(false);
  const [code, setCode] = useState("");
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    if (codePath) {
      fetchCodeByPath();
    } else if (metadata) {
      fetchCodeByMetadata(metadata);
    }
  }, []);

  const fetchCodeByPath = () => {
    DataService.getCodeByPath({ filePath: codePath })
      .then(code => {
        setCode(code);
        setShowModal(true);
      })
      .finally(() => {
        setShowModal(true);
        setTimeout(() => Prism.highlightAll(), 100);
      });
  };

  const fetchCodeByMetadata = () => {
    DataService.getCodeByMetadata(metadata)
      .then(code => {
        setCode(code);
        setShowModal(true);
      })
      .finally(() => {
        setShowModal(true);
        setTimeout(() => Prism.highlightAll(), 100);
      });
  };

  const updateCodeByFilePath = () => {
    return DataService.updateCodeByFilePath({ filePath: codePath, code: code });
  };

  const updateCodeByProblemPath = () => {
    return DataService.updateCodeByProblemPath(metadata, { code: code });
  };

  const updateAndCloseModal = () => {
    if (codePath) {
      updateCodeByFilePath().then(resp => closeModal());
    } else if (metadata) {
      updateCodeByProblemPath().then(resp => closeModal());
    }
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
        <Editor
          value={code}
          highlight={code => highlight(code, languages.js)}
          onValueChange={code => {
            setCode(code);
            setIsCodeUpdated(true);
          }}
          padding={10}
          tabSize={4}
          style={{
            fontFamily: '"Fira code", "Fira Mono", monospace',
            fontSize: 13
          }}
        />
      </Modal.Body>
      <Modal.Footer>
        {isCodeUpdated && (
          <Button
            size="sm"
            variant="outline-success"
            onClick={() => updateAndCloseModal()}
          >
            Save Changes and Close
          </Button>
        )}
        <Button size="sm" variant="outline-danger" onClick={() => closeModal()}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
