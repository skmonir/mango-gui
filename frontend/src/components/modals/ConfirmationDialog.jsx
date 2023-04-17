import React, { useState } from "react";
import PropTypes from "prop-types";
import { Modal, Button, Row, Col, Form } from "react-bootstrap";
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
  flag = {
    show: false,
    label: "",
  },
  show,
  proceed,
}) => {
  const [form, setForm] = useState({
    ok: false,
    flag: false,
  });

  return (
    <Modal
      size="md"
      aria-labelledby="contained-modal-title-vcenter"
      show={show}
      keyboard={false}
      centered
    >
      <Modal.Body style={{ paddingBottom: "2px", paddingTop: "5px" }}>
        <Row className="d-flex mb-2" style={{ fontSize: 22, fontWeight: 650 }}>
          <Col xs={12}>{title}</Col>
        </Row>
        <Row className="mb-3">
          <Col xs={12}>{message}</Col>
        </Row>
        {flag.show && (
          <Row>
            <Col xs={12}>
              <Form.Check
                type="checkbox"
                label={flag.label}
                onChange={(e) => {
                  setForm({ ...form, flag: e.currentTarget.checked });
                }}
              />
            </Col>
          </Row>
        )}
      </Modal.Body>
      <Modal.Footer style={{ paddingBottom: "2px", paddingTop: "2px" }}>
        <Button
          size="sm"
          variant={cancelButton.variant}
          onClick={() => proceed({ ...form, ok: false })}
        >
          {cancelButton.label}
        </Button>
        <Button
          size="sm"
          variant={okButton.variant}
          onClick={() => {
            proceed({ ...form, ok: true });
          }}
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
  flag: PropTypes.object,
  show: PropTypes.bool,
  proceed: PropTypes.func, // called when ok button is clicked.
  enableEscape: PropTypes.bool,
};

export function confirmDialog(options = {}) {
  return createConfirmation(confirmable(ConfirmationDialog))({
    ...options,
  });
}
