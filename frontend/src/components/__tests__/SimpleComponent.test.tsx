import { render, screen } from "@testing-library/react";

// Simple component without hooks
function SimpleComponent({ title }: { title: string }) {
  return <h1>{title}</h1>;
}

describe("SimpleComponent", () => {
  it("should render the title", () => {
    render(<SimpleComponent title="Test Title" />);
    expect(screen.getByText("Test Title")).toBeInTheDocument();
  });
});
