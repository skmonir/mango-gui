import { Col, Row, Spinner } from "react-bootstrap";

export default function Loading() {
  return (
    <Row className="no-gutters">
      <Col xs="12" className="d-flex justify-content-center">
        <Spinner
          variant="success"
          style={{ width: "7rem", height: "7rem", margin: "5px" }}
        />
      </Col>
    </Row>
  );
}
