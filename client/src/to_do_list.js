import React, {Component} from "react";
import axios from "axios";
import {Card, Header, Form, Input, Icon, Button} from "semantic-ui-react";

let endPoint = "http://localhost:9000"

class ToDoList extends Component{
    constructor(props) {
        super(props);

        this.state ={
            task:"",
            items:[],
        };
    }
    componentDidMount(){
        this.getTask();
    }

    onChange = (event) => {
        this.setState({
           [event.target.name] : event.target.value,
        });
    };

    onSubmit = ()=>{
        let {task} = this.state;

        axios.post(endPoint + "/task/v1/create", {task,}, {
            headers:{
                "Content-Type":"application/x-www-from-urlencoded",
            },
        }).then((res)=>{
            this.getTask()
            this.setState({
                task:"",
            })
            console.log(res)
        });
    }

    getTask = ()=>{
        axios.get(endPoint + "/task/v1/list").then((res) =>{
            if (res.data) {
                this.setState({
                    items: res.data.map((item)=>{
                        let color = "yellow";
                        let style = {
                            wordWrap: "break-word",
                        };

                        if(item.status) {
                            color="green";
                            style["textDecorationLine"] = "line-through";
                        }

                        return(
                            <Card key={item._id} color={color} fluid className="rough">
                                <Card.Content>
                                    <Card.Header textAlign="left">
                                        <div style={style}>{item.task}</div>
                                    </Card.Header>

                                    <Card.Meta textAlign="right">
                                        <Icon
                                            name="done"
                                            color="green"
                                            onClick={() => this.doneTask(item._id)}
                                        />
                                        <span style={{paddingRight: 10}}>Done</span>

                                        <Icon
                                            name="undone"
                                            color="blue"
                                            onClick={() => this.undoneTask(item._id)}
                                        />
                                        <span style={{paddingRight: 10}}>Undone</span>

                                        <Icon
                                            name="delete"
                                            color="red"
                                            onClick={() => this.deleteTask(item._id)}
                                        />
                                        <span style={{paddingRight: 10}}>Delete</span>
                                    </Card.Meta>
                                </Card.Content>
                            </Card>
                        );
                    }),
                });
            }else{
                this.setState({
                    item:[],
                });
            }
        });
    };

    doneTask = (id)=> {
        axios.put(endPoint + "/task/v1/done/" + id, {
            headers:{
                "Content-Type":"application/x-www-from-urlencoded",
            },
        }).then((res)=>{
           console.log(res)
           this.getTask()
        });
    }

    undoneTask = (id)=> {
        axios.put(endPoint + "/task/v1/undone/" + id, {
            headers:{
                "Content-Type":"application/x-www-from-urlencoded",
            },
        }).then((res)=>{
            console.log(res)
            this.getTask()
        });
    }

    deleteTask = (id)=> {
        axios.delete(endPoint + "/task/v1/delete/" + id, {
            headers:{
                "Content-Type":"application/x-www-from-urlencoded",
            },
        }).then((res)=>{
            console.log(res)
            this.getTask()
        });
    }

    render(){
        return(
        <div>
            <div className="row">
                <Header className="header" as="h2" color="yellow">
                  TO DO LIST
                </Header>
            </div>
            <div className="row">
                <Form onSubmit={this.onSubmit}>
                    <Input
                    type="text"
                    name="task"
                    onChange={this.onChange}
                    value={this.state.task}
                    fluid
                    placeholder="Create task"
                    />
                    {<Button>Create Task</Button>}
                </Form>
            </div>
            <div className="row">
                <Card.Group>{this.state.items}</Card.Group>
            </div>

        </div>
        );
    };
}

export default ToDoList;
