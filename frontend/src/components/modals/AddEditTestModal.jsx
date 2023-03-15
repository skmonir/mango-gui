import { useEffect, useState } from "react";
import { Alert, Button, Col, Modal, Row } from "react-bootstrap";
import DataService from "../../services/DataService.js";
import Form from "react-bootstrap/Form";
import ShowToast from "../Toast/ShowToast.jsx";

export default function AddEditTestModal({
  metadata,
  testcaseFilePath,
  closeAddEditTestModal
}) {
  const [eventType, setEventType] = useState("");
  const [inputOutputObj, setInputOutputObj] = useState({
    input: "",
    output: ""
  });
  const [showModal, setShowModal] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [showToast, setShowToast] = useState(false);

  useEffect(() => {
    console.log(testcaseFilePath);
    if (testcaseFilePath) {
      setEventType("Update");
      fetchTestcaseByFilePath();
    } else {
      setEventType("Save");
      setShowModal(true);
    }
  }, []);

  const fetchTestcaseByFilePath = () => {
    DataService.getTestcaseByFilePath(testcaseFilePath)
      .then(testcase => {
        setShowModal(true);
        setTimeout(() => {
          setInputOutputObj({
            input: testcase.input,
            output: testcase.output
          });
        }, 0);
      })
      .catch(e => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the testcase!"
        );
      })
      .finally(() => setShowModal(true));
  };

  const saveTriggered = isCloseAfterSave => {
    const data = metadata.split("/");
    let req = {
      platform: data[0],
      contestId: data[1],
      label: data[2],
      input: inputOutputObj.input,
      output: inputOutputObj.output
    };
    if (eventType === "Save") {
      DataService.addCustomTest(req)
        .then(resp => {
          if (resp.message === "success") {
            showToastMessage("Success", "Saved custom testcase successfully!");
            if (isCloseAfterSave) {
              setTimeout(() => closeModal(), 1000);
            } else {
              setInputOutputObj({
                input: "",
                output: ""
              });
            }
          } else {
            showToastMessage("Error", "Error from server!");
          }
        })
        .catch(() => showToastMessage("Error", "Error from server!"));
    } else if (eventType === "Update") {
      DataService.updateCustomTest({
        ...req,
        inputFilePath: testcaseFilePath.inputFilePath,
        outputFilePath: testcaseFilePath.outputFilePath
      })
        .then(resp => {
          if (resp.message === "success") {
            showToastMessage(
              "Success",
              `${eventType}d custom testcase successfully!`
            );
            setTimeout(() => closeModal(), 1000);
          } else {
            showToastMessage("Error", "Error from server!");
          }
        })
        .catch(() => showToastMessage("Error", "Error from server!"));
    }
  };

  const closeModal = () => {
    setShowModal(false);
    setTimeout(() => closeAddEditTestModal(), 500);
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
        fullscreen={true}
      >
        <Modal.Body style={{ height: "80vh", overflowY: "auto" }}>
          <Row>
            <Col xs={6}>
              <Form>
                <Form.Label>
                  <strong>INPUT</strong>
                </Form.Label>
                <pre>
                  <Form.Control
                    value={inputOutputObj.input}
                    onChange={e =>
                      setInputOutputObj({
                        ...inputOutputObj,
                        input: e.target.value
                      })
                    }
                    autoCorrect="off"
                    autoComplete="off"
                    autoCapitalize="none"
                    as="textarea"
                    aria-label="With textarea"
                    rows={23}
                  />
                </pre>
              </Form>
            </Col>
            <Col xs={6}>
              <Form>
                <Form.Label>
                  <strong>EXPECTED OUTPUT</strong>
                </Form.Label>
                <pre>
                  <Form.Control
                    value={inputOutputObj.output}
                    onChange={e =>
                      setInputOutputObj({
                        ...inputOutputObj,
                        output: e.target.value
                      })
                    }
                    autoCorrect="off"
                    autoComplete="off"
                    autoCapitalize="none"
                    as="textarea"
                    aria-label="With textarea"
                    rows={23}
                  />
                </pre>
              </Form>
            </Col>
          </Row>
        </Modal.Body>
        <Modal.Footer>
          <Button
            size="sm"
            variant="outline-primary"
            disabled={
              !inputOutputObj || !inputOutputObj.input || !inputOutputObj.output
            }
            onClick={() => saveTriggered(true)}
          >
            {`${eventType} and Close`}
          </Button>
          {eventType === "Save" && (
            <Button
              size="sm"
              variant="outline-success"
              disabled={
                !inputOutputObj ||
                !inputOutputObj.input ||
                !inputOutputObj.output
              }
              onClick={() => saveTriggered(false)}
            >
              {`Save and Add Another One`}
            </Button>
          )}
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
