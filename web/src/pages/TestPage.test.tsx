import '@testing-library/dom'
import { render, screen } from "@testing-library/react";
import { TestPage } from "./TestPage";

describe("TestPage", () => {
  it('should render', () => {
    render(<TestPage />);
    screen.getByText('Hello');
  })
});
