import {
  Alert,
  Button,
  Col,
  InputGroup,
  Modal,
  Row,
  Spinner
} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEye,
  faEyeSlash,
  faRightToBracket
} from "@fortawesome/free-solid-svg-icons";
import DataService from "../../services/DataService.js";
import { useState } from "react";
import Utils from "../../Utils.js";

export default function LoginModal({ setShowLoginModal }) {
  const [form, setForm] = useState({
    platform: "codeforces",
    handleOrEmail: "",
    password: ""
  });
  const [showModal, setShowModal] = useState(true);
  const [showPass, setShowPass] = useState(false);
  const [loginInProgress, setLoginInProgress] = useState(false);
  const [alert, setAlert] = useState({
    message: "",
    variant: ""
  });

  const validate = () => {
    let ok = true;
    for (const [_, value] of Object.entries(form)) {
      ok = ok && !Utils.isStrNullOrEmpty(value);
    }
    if (!ok) {
      setAlert({
        message: "No field can be empty",
        variant: "danger"
      });
    }
    return ok;
  };

  const doLogin = () => {
    console.log(form);
    if (!validate()) {
      return;
    }
    setLoginInProgress(true);
    DataService.login(form)
      .then(resp => {
        console.log(resp);
        setAlert({
          message: `Hi, ${resp.handle}!`,
          variant: "success"
        });
        setTimeout(() => closeModal(), 1000);
      })
      .catch(error => {
        setAlert({
          message: error.response.data.message,
          variant: "danger"
        });
      })
      .finally(() => setLoginInProgress(false));
  };

  const closeModal = () => {
    setShowModal(false);
    setTimeout(() => setShowLoginModal(false), 500);
  };

  return (
    <Modal
      show={showModal}
      size="sm"
      aria-labelledby="contained-modal-title-vcenter"
      centered
    >
      <Modal.Header style={{ paddingBottom: "5px", paddingTop: "5px" }}>
        <strong>Login</strong>
      </Modal.Header>
      <Modal.Body style={{ paddingBottom: "2px", paddingTop: "5px" }}>
        <Row>
          <Col xs={12}>
            <Form.Group className="mb-3">
              <Form.Label>Platform</Form.Label>
              <Form.Select size="sm" aria-label="Default select example">
                <option value="codeforces">Codeforces</option>
              </Form.Select>
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <Form.Group className="mb-3">
              <Form.Label>Handle or Email</Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                value={form.handleOrEmail}
                onChange={e =>
                  setForm({ ...form, handleOrEmail: e.target.value })
                }
              />
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <Form.Group className="mb-3">
              <Form.Label>Password</Form.Label>
              <InputGroup className="mb-3">
                <Form.Control
                  type={showPass ? "text" : "password"}
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  value={form.password}
                  onChange={e => setForm({ ...form, password: e.target.value })}
                />
                <Button
                  size="sm"
                  variant="secondary"
                  onClick={() => setShowPass(!showPass)}
                >
                  <FontAwesomeIcon icon={showPass ? faEyeSlash : faEye} />
                </Button>
              </InputGroup>
            </Form.Group>
          </Col>
        </Row>
        {alert.message && (
          <Row>
            <Col>
              <Alert variant={alert.variant} className="text-center p-1 mb-2">
                {alert.message}
              </Alert>
            </Col>
          </Row>
        )}
      </Modal.Body>
      <Modal.Footer style={{ paddingBottom: "5px", paddingTop: "5px" }}>
        <Button size="sm" variant="outline-secondary" onClick={closeModal}>
          Cancel
        </Button>
        <Button size="sm" variant="outline-success" onClick={doLogin}>
          {loginInProgress ? (
            <Spinner
              as="span"
              animation="border"
              size="sm"
              role="status"
              aria-hidden="true"
            />
          ) : (
            <FontAwesomeIcon icon={faRightToBracket} />
          )}
          {" Login"}
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
