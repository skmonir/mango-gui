import {Col, Row, Spinner} from "react-bootstrap";

export default function Loading() {

    return (
        <Row className="no-gutters">
            <Col sm="12" md={{ size: 6, offset: 5 }}>
                <Spinner
                    color="primary"
                    style={{ width: "5rem", height: "5rem", margin: "5px" }}
                />
            </Col>
        </Row>
    );
}