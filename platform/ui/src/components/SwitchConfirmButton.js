import React from "react";
import {Button, OverlayTrigger, Popover} from "react-bootstrap";

export default function SwitchConfirmButton({onClick, enabled, children, allowSwitchContainer}) {
  const [startUpgrade, setStartUpgrade] = React.useState();

  const handleClick = React.useCallback((e) => {
    e.preventDefault();

    setStartUpgrade(false);
    onClick();
  }, [onClick]);

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
              onClick={(e) => handleClick(e)}
            >
              确认
            </Button>
          </div>
          <div className="col-12">
            <Button
              variant="primary"
              onClick={() => setStartUpgrade(false)}
            >
              取消
            </Button>
          </div>
        </div>
      </Popover.Body>
    </Popover>
  );

  if (!enabled) return <></>;
  return (<>
    <OverlayTrigger trigger="click" placement="right" overlay={popover} show={startUpgrade}>
      <Button
        variant='warning'
        disabled={!allowSwitchContainer}
        onClick={() => setStartUpgrade(true)}
      >
        切换
      </Button>
    </OverlayTrigger>
  </>);
}

