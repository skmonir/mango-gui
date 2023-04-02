import { useEffect, useState } from "react";
import { Alert, Button, Col, Modal, Row } from "react-bootstrap";
import DataService from "../../services/DataService.js";
import Form from "react-bootstrap/Form";
import ShowToast from "../Toast/ShowToast.jsx";
import { faPlus, faSave } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function AddEditCustomTestModal({
  event,
  metadata,
  testcaseFilePath,
  closeAddEditTestModal
}) {
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
    if (event === "Update" || event === "Clone") {
      fetchTestcaseByFilePath();
    } else if (event === "Add") {
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
    if (event !== "Update") {
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
    } else if (event === "Update") {
      DataService.updateCustomTest({
        ...req,
        inputFilePath: testcaseFilePath.inputFilePath,
        outputFilePath: testcaseFilePath.outputFilePath
      })
        .then(resp => {
          if (resp.message === "success") {
            showToastMessage(
              "Success",
              "Updated custom testcase successfully!"
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
      >
        <Modal.Header
          style={{ paddingBottom: "5px", paddingTop: "5px" }}
          closeButton
        >
          <strong>{`${event} Custom Test`}</strong>
        </Modal.Header>
        <Modal.Body
          style={{
            overflowY: "auto",
            paddingBottom: "1px",
            paddingTop: "1px"
          }}
        >
          <Row>
            <Col xs={6}>
              <Form>
                <Form.Label>
                  <strong>INPUT</strong>
                </Form.Label>
                <pre>
                  <Form.Control
                    style={{ fontSize: 14 }}
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
                    rows={15}
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
                    style={{ fontSize: 14 }}
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
                    rows={15}
                  />
                </pre>
              </Form>
            </Col>
          </Row>
        </Modal.Body>
        <Modal.Footer style={{ paddingBottom: "2px", paddingTop: "2px" }}>
          <Button
            size="sm"
            variant="outline-primary"
            disabled={
              !inputOutputObj || !inputOutputObj.input || !inputOutputObj.output
            }
            onClick={() => saveTriggered(true)}
          >
            <FontAwesomeIcon icon={faSave} />{" "}
            {(event !== "Update" ? "Save" : "Update") + ` and Close`}
          </Button>
          {event === "Add" && (
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
              <FontAwesomeIcon icon={faSave} />{" "}
              <FontAwesomeIcon icon={faPlus} /> Save and Add Another One
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
