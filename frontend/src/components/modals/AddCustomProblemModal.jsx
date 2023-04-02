import { Alert, Button, Col, Modal, Row } from "react-bootstrap";
import ShowToast from "../Toast/ShowToast.jsx";
import { useEffect, useState } from "react";
import Form from "react-bootstrap/Form";
import Utils from "../../Utils.js";
import DataService from "../../services/DataService.js";
import { faSave } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function AddCustomProblemModal({
  closeAddCustomProblemModal,
  insertCustomProblemIntoList
}) {
  const [showModal, setShowModal] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [showToast, setShowToast] = useState(false);
  const [customProblem, setCustomProblem] = useState({
    platform: "custom",
    contestId: "",
    label: "",
    name: "",
    timeLimit: 2,
    memoryLimit: 512,
    url: "",
    status: ""
  });

  useEffect(() => {
    setShowModal(true);
  }, []);

  const closeModal = () => {
    setShowModal(false);
    setTimeout(() => closeAddCustomProblemModal(), 500);
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
  };

  const validate = () => {
    let errorMessage = "";
    const strKeys = ["platform", "contestId", "label", "name"];
    strKeys.forEach(key => {
      if (Utils.isStrNullOrEmpty(customProblem[key])) {
        errorMessage = "No field can be empty\n";
      }
    });
    if (!Utils.isValidNum(customProblem.timeLimit, 1, 10)) {
      errorMessage += "Time Limit should be an integer\n";
    }
    if (!Utils.isValidNum(customProblem.memoryLimit, 128, 2048)) {
      errorMessage += "Memory Limit should be a integer\n";
    }
    if (errorMessage !== "") {
      showToastMessage("Error", errorMessage);
      return false;
    }
    return true;
  };

  const saveAndCloseModal = () => {
    if (validate()) {
      DataService.addCustomProblem(customProblem)
        .then(data => {
          insertCustomProblemIntoList(data);
          closeModal();
        })
        .catch(e => {
          showToastMessage(
            "Error",
            "Oops! Something went wrong while saving the problem!"
          );
        });
    }
  };

  return (
    <div>
      <Modal
        show={showModal}
        onHide={closeModal}
        size="lg"
        aria-labelledby="contained-modal-title-vcenter"
        centered
      >
        <Modal.Header
          style={{ paddingBottom: "5px", paddingTop: "5px" }}
          closeButton
        >
          <strong>Add Custom Problem</strong>
        </Modal.Header>
        <Modal.Body
          style={{
            height: "55vh",
            overflowY: "auto",
            paddingBottom: "1px",
            paddingTop: "1px"
          }}
        >
          <Row>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Platform</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Example. 1502, abc123"
                  value={customProblem.platform}
                  disabled={true}
                  onChange={e =>
                    setCustomProblem({
                      ...customProblem,
                      platform: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Contest ID</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Example. 1502, abc123"
                  value={customProblem.contestId}
                  onChange={e =>
                    setCustomProblem({
                      ...customProblem,
                      contestId: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Time Limit [Sec]</strong>
                </Form.Label>
                <Form.Control
                  type="number"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  value={customProblem.timeLimit}
                  onChange={e =>
                    setCustomProblem({
                      ...customProblem,
                      timeLimit: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Memory Limit [MB]</strong>
                </Form.Label>
                <Form.Control
                  type="number"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  value={customProblem.memoryLimit}
                  onChange={e =>
                    setCustomProblem({
                      ...customProblem,
                      memoryLimit: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
          </Row>
          <Row>
            <Col xs={2}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Label</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  value={customProblem.label}
                  onChange={e =>
                    setCustomProblem({
                      ...customProblem,
                      label: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col xs={10}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Problem Name</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  value={customProblem.name}
                  onChange={e =>
                    setCustomProblem({
                      ...customProblem,
                      name: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
          </Row>
          <Row>
            <Alert variant="info">
              <pre>
                {`*** Follow any of the ways below to generate testcases for this problem. \n1. Use Input Generator \n2. Load problem in Tester and add custom test\n\n`}
                {`*** Custom platform URL will be like 'custom/:contest_id/:problem_label'\nSo use like 'custom/1234/A', 'custom/1234/D1' in Tester or Input/Output Generator`}
              </pre>
            </Alert>
          </Row>
        </Modal.Body>
        <Modal.Footer style={{ paddingBottom: "2px", paddingTop: "2px" }}>
          <Button
            size="sm"
            variant="outline-success"
            onClick={() => saveAndCloseModal()}
          >
            <FontAwesomeIcon icon={faSave} /> Save Problem
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
