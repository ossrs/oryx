//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import SrsQRCode from "./SrsQRCode";

test('renders qrCode component', () => {
  render(<SrsQRCode url='http://localhost' />);
  const elem = screen.getByTestId(/qrCode/i);
  expect(elem).toBeInTheDocument();
});

test('renders qrCode without url', () => {
  render(<SrsQRCode />);
  const elem = screen.queryByTestId(/qrCode/i);
  expect(elem).toBeNull();
});

