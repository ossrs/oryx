// See https://stackoverflow.com/a/56557915/17679565
import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
import SrsQRCode from "./SrsQRCode";

test('renders qrCode component', () => {
  render(<SrsQRCode url='http://localhost' />);
  const elem = screen.getByTestId(/qrCode/i);
  expect(elem).toBeInTheDocument();
});

