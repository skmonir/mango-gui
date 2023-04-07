import { Col, Modal, Row } from "react-bootstrap";
import Loading from "../misc/Loading.jsx";

export default function InitializerModal({ showModal, initMessage }) {
  return (
    <Modal show={showModal} size="sm" centered>
      <Modal.Body style={{ paddingBottom: "2px", paddingTop: "5px" }}>
        <Row>
          <Col xs={12}>
            <Loading />
          </Col>
        </Row>
        <Row
          className="d-flex text-center"
          style={{ fontSize: 22, color: "darkcyan", fontWeight: 500 }}
        >
          <pre>{initMessage}</pre>
        </Row>
        <Row className="d-flex text-center">
          <pre>Please wait a moment</pre>
        </Row>
      </Modal.Body>
    </Modal>
  );
}
