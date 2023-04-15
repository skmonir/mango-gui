import React from "react";
import PropTypes from "prop-types";
import { Modal, Button, Row, Col } from "react-bootstrap";
import { confirmable, createConfirmation } from "react-confirm";

const ConfirmationDialog = ({
  title = "Confirmation!",
  message = "Are you sure?",
  okButton = {
    label: "OK",
    variant: "success",
  },
  cancelButton = {
    label: "Cancel",
    variant: "secondary",
  },
  show,
  proceed,
}) => {
  return (
    <Modal
      size="md"
      aria-labelledby="contained-modal-title-vcenter"
      show={show}
      onHide={() => proceed(false)}
      centered
    >
      <Modal.Body style={{ paddingBottom: "2px", paddingTop: "5px" }}>
        <Row className="d-flex" style={{ fontSize: 22, fontWeight: 650 }}>
          <Col xs={12}>
            <pre>{title}</pre>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <pre style={{ whiteSpace: "pre-wrap" }}>{message}</pre>
          </Col>
        </Row>
      </Modal.Body>
      <Modal.Footer style={{ paddingBottom: "2px", paddingTop: "2px" }}>
        <Button
          size="sm"
          variant={cancelButton.variant}
          onClick={() => proceed(false)}
        >
          {cancelButton.label}
        </Button>
        <Button
          size="sm"
          variant={okButton.variant}
          onClick={() => proceed(true)}
        >
          {okButton.label}
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

ConfirmationDialog.propTypes = {
  title: PropTypes.string,
  message: PropTypes.string,
  okButton: PropTypes.object,
  cancelButton: PropTypes.object,
  show: PropTypes.bool,
  proceed: PropTypes.func, // called when ok button is clicked.
  enableEscape: PropTypes.bool,
};

export function confirmDialog(options = {}) {
  return createConfirmation(confirmable(ConfirmationDialog))({
    ...options,
  });
}
