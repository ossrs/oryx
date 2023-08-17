//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import {BrowserRouter} from "react-router-dom";
import {SrsErrorBoundary} from "./SrsErrorBoundary";
import {Errors} from "../utils";

// See https://github.com/facebook/react/issues/11098#issuecomment-523977830
function disableConsoleError(pfn) {
  const spy = jest.spyOn(console, 'error');
  spy.mockImplementation(() => {});

  try {
    pfn();
  } finally {
    spy.mockRestore();
  }
}

test('renders with error', () => {
  const TestError = () => {
    throw new Error();
  }

  disableConsoleError(() => {
    render(<SrsErrorBoundary><TestError/></SrsErrorBoundary>);
    const elem = screen.getByText(/Name: Error/i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with error desc', () => {
  const TestError = () => {
    throw new Error('SomethingFailByUTest');
  }

  disableConsoleError(() => {
    render(<SrsErrorBoundary><TestError /></SrsErrorBoundary>);
    const elem = screen.getByText(/Message: SomethingFailByUTest/i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with auth', () => {
  const TestError = () => {
    throw {response: {data: {code: Errors.auth}}};
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Token过期/i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with code', () => {
  const TestError = () => {
    throw {response: {data: {code: 100}}};
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Code: 100/i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with status', () => {
  const TestError = () => {
    throw {response: {status: 500}};
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Status: 500/i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with object', () => {
  const TestError = () => {
    throw {key: 100};
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Object: /i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with array', () => {
  const TestError = () => {
    throw [1, 'a', {}];
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Object: /i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with funciton', () => {
  const TestError = () => {
    throw ()=>{};
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Function: /i);
    expect(elem).toBeInTheDocument();
  });
});

test('renders with string', () => {
  const TestError = () => {
    throw 'Hello World!';
  }

  disableConsoleError(() => {
    render(<BrowserRouter><SrsErrorBoundary><TestError /></SrsErrorBoundary></BrowserRouter>);
    const elem = screen.getByText(/Hello World!/i);
    expect(elem).toBeInTheDocument();
  });
});

