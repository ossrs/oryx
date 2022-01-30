import React from "react";
import {Button, Spinner, OverlayTrigger, Popover} from "react-bootstrap";

export default function PopoverConfirmButton({upgrading, handleClick, operator, text, children}) {
  const [showUpgrading, setShowUpgrading] = React.useState();

  const onHandleClick = () => {
    setShowUpgrading(false);
    handleClick();
  };

  const filterByOperator = (msg) => {
    return operator ? msg.replaceAll('升级', operator) : msg;
  };

  const popover = (
    <Popover id="popover-basic">
      <Popover.Header as="h3">Confirm</Popover.Header>
      <Popover.Body>
        <div>
          {children}
        </div>
        <div className='row row-cols-lg-auto g-3 align-items-center'>
          <div className="col-12">
            <Button
              variant="danger"
              disabled={upgrading}
              onClick={!upgrading ? onHandleClick : null}
            >
              确认
            </Button>
          </div>
          <div className="col-12">
            <Button
              variant="primary"
              onClick={() => setShowUpgrading(false)}
            >
              取消
            </Button>
          </div>
        </div>
      </Popover.Body>
    </Popover>
  );

  return (
    <div className='row row-cols-lg-auto g-3 align-items-center' style={{display: 'inline-block'}}>
      <div className="col-12">
        <OverlayTrigger trigger="click" placement="right" overlay={popover} show={showUpgrading}>
          <Button variant="primary" onClick={() => setShowUpgrading(!showUpgrading)}>
            {upgrading ? filterByOperator('升级中...') : text}
          </Button>
        </OverlayTrigger> &nbsp;
        {upgrading && <Spinner animation="border" style={{verticalAlign: 'middle'}} />}
      </div>
    </div>
  );
}

