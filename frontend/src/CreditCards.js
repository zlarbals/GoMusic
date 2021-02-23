import React from "react";
import {
  injectStripe,
  StripeProvider,
  Elements,
  CardElement,
} from "react-stripe-elements";
import cookie from "js-cookie";

const INITIALSTATE = "INITIAL",
  SUCCESSSTATE = "COMPLETE",
  FAILEDSTATE = "FAILED";

class CreditCardForm extends React.Component {
  constructor(props) {
    super(props);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.state = {
      status: INITIALSTATE,
      useExisting: false,
    };
  }

  renderCreditCardInformation() {
    const user = cookie.getJSON("user");

    const style = {
      base: {
        fontSize: "20px",
        color: "#495057",
        fontFamily:
          'apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif',
      },
    };

    const usersavedcard = (
      <div>
        <div className="form-row text-center">
          <button
            type="submit"
            className="btn btn-outline-success text-center mx-auto"
            onClick={() => this.setState({ useExisting: true })}
          >
            Use saved card?
          </button>
        </div>
        <hr />
      </div>
    );

    let remembercardcheck = null;
    if (user.loggedin === true) {
      remembercardcheck = (
        <div className="form-row form-check text-center">
          <input
            className="form-check-input"
            type="checkbox"
            value=""
            id="remembercardcheck"
            name="remember"
            onChange={this.handleChange}
          />
          <label className="form-check-label" htmlFor="remembercardcheck">
            Remember Card?
          </label>
        </div>
      );
    }

    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          {user.loggedin ? usersavedcard : null}
          <h5 className="mb-4">Payment Info</h5>
          <div className="form-row">
            <div className="col-lg-12 form-group">
              <label htmlFor="cc-name">Name On Card:</label>
              <input
                id="cc-name"
                name="name"
                className="form-control"
                placeholder="Name on Card"
                onChange={this.handleChange}
                type="text"
              />
            </div>
          </div>
          <div className="form-row">
            <div className="col-lg-12 form-group">
              <label htmlFor="card">Card Information:</label>
              <CardElement id="card" className="form-control" style={style} />
            </div>
          </div>
          {remembercardcheck}
          <hr className="mb-4" />
          <button type="sumbit" className="btn btn-success btn-large">
            {this.props.operation}
          </button>
        </form>
      </div>
    );
  }

  renderSuccess() {
    return (
      <div>
        <h5 className="mb-4 text-success">Request Successfull....</h5>
        <button
          type="submit"
          className="btn btn-success btn-large"
          onClick={() => {
            this.props.toggle();
          }}
        >
          Ok
        </button>
      </div>
    );
  }
  renderFailure() {
    return (
      <div>
        <h5 className="mb-4 text-danger">
          Credit card information invalid, try again or exit
        </h5>
        {this.renderCreditCardInformation()}
      </div>
    );
  }

  async handleSubmit(event) {
    event.preventDefault();
    let id = "";
    console.log(this.state);

    //저장된 카드 사용이 아니라면 스트라이프에 토큰을 요청한다.
    if (!this.state.useExisting) {
      //스트파이프 API를 사용해 토큰을 생성한다.
      console.log("hereherehere");
      let { token } = await this.props.stripe.createToken({
        name: this.state.name,
      });
      if (token == null) {
        console.log("invalid token");
        this.setState({ status: FAILEDSTATE });
        return;
      }
      id = token.id;
    }

    console.log("fetch 직전 : " + this.state.useExisting);
    console.log(this.state);

    //요청을 생성하고 백엔드로 보낸다.
    let response = await fetch("/users/charge", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        token: id,
        customer_id: this.props.user,
        product_id: this.props.productid,
        sell_price: this.props.price,
        rememberCard: this.state.remember !== undefined,
        useExisting: this.state.useExisting,
      }),
    });

    //응답이 ok면 작업 성공
    if (response.ok) {
      console.log("Purchase Complete!");
      this.setState({ status: SUCCESSSTATE, useExisting: false });
    } else {
      this.setState({ status: FAILEDSTATE, useExisting: false });
    }
  }

  handleChange(event) {
    event.preventDefault();
    const name = event.target.name;
    const value = event.target.value;
    this.setState({
      [name]: value,
    });
  }

  render() {
    let body = null;
    switch (this.state.status) {
      case SUCCESSSTATE:
        body = this.renderSuccess();
        break;
      case FAILEDSTATE:
        body = this.renderFailure();
        break;
      default:
        body = this.renderCreditCardInformation();
    }
    return <div>{body}</div>;
  }
}

export default function CreditCardInformation(props) {
  if (!props.show) {
    return null;
  }

  //스트라이프 API를 사용해 CreditCardForm를 추가하면 createToken() 메서드를 호출할 수 있다.
  const CCFormWithStripe = injectStripe(CreditCardForm);
  return (
    <div>
      {props.separator ? <hr /> : null}
      {/*input stripe public key*/}
      <StripeProvider apiKey="pubilc key">
        <Elements>
          <CCFormWithStripe
            user={props.user}
            operation={props.operation}
            productid={props.productid}
            price={props.price}
            toggle={props.toggle}
          />
        </Elements>
      </StripeProvider>
    </div>
  );
}
