import { useEffect, useState } from "react";
import { Button, Modal } from "react-bootstrap";
import DataService from "../../services/DataService.js";
import Prism from "prismjs";
import "prismjs/themes/prism.css";
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import Editor from "react-simple-code-editor";
import ShowToast from "../Toast/ShowToast.jsx";
import {
  faCompress,
  faMaximize,
  faSave
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function ViewCodeModal({
  codePath,
  metadata,
  setShowCodeModal
}) {
  const [isCodeUpdated, setIsCodeUpdated] = useState(false);
  const [code, setCode] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [showToast, setShowToast] = useState(false);
  const [isFullScreen, setIsFullScreen] = useState(false);

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
      .catch(e => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the code!"
        );
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
      .catch(e => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the code!"
        );
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
      updateCodeByFilePath()
        .then(resp => closeModal())
        .catch(e => {
          showToastMessage(
            "Error",
            "Oops! Something went wrong while saving the code!"
          );
        });
    } else if (metadata) {
      updateCodeByProblemPath()
        .then(resp => closeModal())
        .catch(e => {
          showToastMessage(
            "Error",
            "Oops! Something went wrong while saving the code!"
          );
        });
    }
  };

  const closeModal = () => {
    setShowModal(false);
    setTimeout(() => setShowCodeModal(false), 500);
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
  };

  return (
    <div>
      <Modal
        show={showModal}
        onHide={closeModal}
        size="lg"
        aria-labelledby="contained-modal-title-vcenter"
        centered
        fullscreen={isFullScreen}
      >
        <Modal.Header closeButton>
          <strong>Code Editor</strong>
        </Modal.Header>
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
          <Button
            size="sm"
            variant="outline-success"
            onClick={() => setIsFullScreen(!isFullScreen)}
          >
            {isFullScreen ? (
              <FontAwesomeIcon icon={faCompress} />
            ) : (
              <FontAwesomeIcon icon={faMaximize} />
            )}
          </Button>
          <Button
            size="sm"
            variant="outline-success"
            disabled={!isCodeUpdated}
            onClick={() => updateAndCloseModal()}
          >
            <FontAwesomeIcon icon={faSave} /> Save Changes and Close
          </Button>
          <Button
            size="sm"
            variant="outline-danger"
            onClick={() => closeModal()}
          >
            Close
          </Button>
        </Modal.Footer>
      </Modal>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
    </div>
  );
}
