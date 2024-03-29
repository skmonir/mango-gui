import { useEffect, useState } from "react";
import { Button, Modal } from "react-bootstrap";
import DataService from "../../services/DataService.js";
import ShowToast from "../Toast/ShowToast.jsx";
import SplitPane, { Pane } from "react-split-pane";
import CodeEditor from "../misc/CodeEditor.jsx";

export default function ViewCodeModal({
  codePath,
  metadata,
  setShowCodeModal,
  executionResult,
  customElementsOnHeader,
}) {
  const [isCodeUpdated, setIsCodeUpdated] = useState(false);
  const [codeContent, setCodeContent] = useState({
    lang: "",
    code: "",
  });
  const [showModal, setShowModal] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: "",
  });
  const [showToast, setShowToast] = useState(false);

  useEffect(() => {
    if (codePath) {
      fetchCodeByPath();
    } else if (metadata) {
      fetchCodeByMetadata(metadata);
    }
  }, []);

  const fetchCodeByPath = () => {
    DataService.getCodeByPath({ filePath: codePath })
      .then((codeResponse) => {
        setCodeContent(codeResponse);
      })
      .catch((e) => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the code!"
        );
      })
      .finally(() => {
        setShowModal(true);
      });
  };

  const fetchCodeByMetadata = () => {
    DataService.getCodeByMetadata(metadata)
      .then((codeResponse) => {
        setCodeContent(codeResponse);
      })
      .catch((e) => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the code!"
        );
      })
      .finally(() => {
        setShowModal(true);
      });
  };

  const updateCodeByFilePath = () => {
    return DataService.updateCodeByFilePath({
      filePath: codePath,
      code: codeContent.code,
    });
  };

  const updateCodeByProblemPath = () => {
    return DataService.updateCodeByProblemPath(metadata, {
      code: codeContent.code,
    });
  };

  const updateCode = () => {
    if (isCodeUpdated) {
      if (codePath) {
        return updateCodeByFilePath();
      } else if (metadata) {
        return updateCodeByProblemPath();
      }
    } else {
      return Promise.resolve(true);
    }
  };

  const updateAndCloseModal = () => {
    updateCode()
      .then((resp) => closeModal())
      .catch((e) => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while saving the code!"
        );
      });
  };

  const onCodeChange = (code) => {
    setCodeContent({
      ...codeContent,
      code: code,
    });
    setIsCodeUpdated(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setTimeout(() => setShowCodeModal(false), 500);
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message,
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
        fullscreen={true}
        keyboard={false}
      >
        <Modal.Header
          style={{ paddingBottom: "5px", paddingTop: "5px" }}
          closeButton
        >
          <strong>Code Editor</strong>
        </Modal.Header>
        <Modal.Body style={{ paddingBottom: "2px", paddingTop: "2px" }}>
          {!executionResult?.show && (
            <CodeEditor
              codeContent={codeContent}
              onChange={onCodeChange}
              onBlur={updateCode}
              readOnly={{ editor: false, preference: false }}
              customElementsOnHeader={customElementsOnHeader}
            />
          )}
          {executionResult?.show && (
            <SplitPane
              split="vertical"
              defaultSize="40%"
              primary="first"
              className="mt-0"
              style={{ height: "99%" }}
            >
              <div
                style={{
                  marginRight: "5px",
                  maxHeight: "95vh",
                  overflowY: "auto",
                }}
              >
                {executionResult.component}
              </div>
              <div style={{ marginLeft: "5px" }}>
                <CodeEditor
                  codeContent={codeContent}
                  onChange={onCodeChange}
                  onBlur={updateCode}
                  readOnly={{ editor: false, preference: false }}
                  customElementsOnHeader={customElementsOnHeader}
                />
              </div>
            </SplitPane>
          )}
        </Modal.Body>
      </Modal>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
    </div>
  );
}
