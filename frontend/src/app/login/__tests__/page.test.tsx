import LoginPage from "@/app/login/page";
import { useKeycloak } from "@react-keycloak/web";
import { render, screen } from "@testing-library/react";
import { useRouter } from "next/navigation";
import React from "react";

// Mock the hooks
jest.mock("@react-keycloak/web", () => ({
  useKeycloak: jest.fn(),
}));
jest.mock("next/navigation", () => ({
  useRouter: jest.fn(),
}));

const mockUseKeycloak = useKeycloak as jest.MockedFunction<typeof useKeycloak>;
const mockUseRouter = useRouter as jest.MockedFunction<typeof useRouter>;

describe("LoginPage", () => {
  let useEffectSpy: jest.SpyInstance;

  beforeEach(() => {
    jest.clearAllMocks();
    useEffectSpy = jest
      .spyOn(React, "useEffect")
      .mockImplementation((f) => f());
    mockUseRouter.mockReturnValue({
      push: jest.fn(),
      replace: jest.fn(),
      back: jest.fn(),
      forward: jest.fn(),
      refresh: jest.fn(),
      prefetch: jest.fn(),
    });
  });

  afterEach(() => {
    useEffectSpy.mockRestore();
  });

  it("renders the main heading when not authenticated", () => {
    mockUseKeycloak.mockReturnValue({
      keycloak: {
        authenticated: false,
        login: jest.fn(),
        logout: jest.fn(),
        register: jest.fn(),
        accountManagement: jest.fn(),
        createLoginUrl: jest.fn(),
        createLogoutUrl: jest.fn(),
        createRegisterUrl: jest.fn(),
        createAccountUrl: jest.fn(),
        isTokenExpired: jest.fn(),
        updateToken: jest.fn(),
        clearToken: jest.fn(),
        hasRealmRole: jest.fn(),
        hasResourceRole: jest.fn(),
        loadUserProfile: jest.fn(),
        loadUserInfo: jest.fn(),
        init: jest.fn(),
        didInitialize: true,
      },
      initialized: true,
    });

    render(<LoginPage />);

    const heading = screen.getByRole("heading", {
      name: /CloudGate SSO/i,
    });
    expect(heading).toBeInTheDocument();
  });

  it("shows initializing message when not initialized", () => {
    mockUseKeycloak.mockReturnValue({
      keycloak: null as any,
      initialized: false,
    });

    render(<LoginPage />);

    expect(screen.getByText(/Initializing.../i)).toBeInTheDocument();
  });

  it("shows redirecting message when authenticated", () => {
    mockUseKeycloak.mockReturnValue({
      keycloak: {
        authenticated: true,
        login: jest.fn(),
        logout: jest.fn(),
        register: jest.fn(),
        accountManagement: jest.fn(),
        createLoginUrl: jest.fn(),
        createLogoutUrl: jest.fn(),
        createRegisterUrl: jest.fn(),
        createAccountUrl: jest.fn(),
        isTokenExpired: jest.fn(),
        updateToken: jest.fn(),
        clearToken: jest.fn(),
        hasRealmRole: jest.fn(),
        hasResourceRole: jest.fn(),
        loadUserProfile: jest.fn(),
        loadUserInfo: jest.fn(),
        init: jest.fn(),
        didInitialize: true,
      },
      initialized: true,
    });

    render(<LoginPage />);

    expect(
      screen.getByText(/Redirecting to dashboard.../i)
    ).toBeInTheDocument();
  });
});
