import { Alert, Button, Col, Modal, Row } from "react-bootstrap";
import ShowToast from "../Toast/ShowToast.jsx";
import { useEffect, useState } from "react";
import Form from "react-bootstrap/Form";

export default function AddCustomProblemModal({ closeAddCustomProblemModal }) {
  const [showModal, setShowModal] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [showToast, setShowToast] = useState(false);

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

  const saveAndCloseModal = () => {};

  return (
    <div>
      <Modal
        show={showModal}
        onHide={closeModal}
        size="lg"
        aria-labelledby="contained-modal-title-vcenter"
        centered
      >
        <Modal.Header />
        <Modal.Body style={{ height: "35vh", overflowY: "auto" }}>
          <Row>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Platform</strong>
                </Form.Label>
                <Form.Select size="sm" aria-label="Default select example">
                  <option value="codeforces">Codeforces</option>
                  <option value="atcoder">Atcoder</option>
                  <option value="custom">Custom</option>
                </Form.Select>
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
                />
              </Form.Group>
            </Col>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Time Limit [Sec]</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                />
              </Form.Group>
            </Col>
            <Col xs={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Memory Limit [MB]</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
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
                />
              </Form.Group>
            </Col>
          </Row>
          <Row>
            <Alert variant="info" className="text-center">
              To add testcase, use Input Generator or load problem in Tester and
              add custom test.
            </Alert>
          </Row>
        </Modal.Body>
        <Modal.Footer>
          <Button
            size="sm"
            variant="outline-success"
            onClick={() => saveAndCloseModal()}
          >
            Save Problem
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
