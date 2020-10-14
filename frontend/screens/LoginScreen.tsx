import React from 'react';
import { StyleSheet, Text,TouchableOpacity,TouchableHighlight,TextInput,View } from 'react-native';
import { string } from 'prop-types';
import { NavigationScreenProp } from 'react-navigation';
import { withNavigation } from 'react-navigation';
import Navigation from '../navigation';


export default class App extends React.Component {
  state={
    username:"",
    password:""
  }
  



  handleLoginOnPress = () => {
    var username = this.state['username'];
    var password = this.state['password']
    var payload = {
      "username": username,
      "password": password
    }
    fetch('http://localhost:8000/signin', {
      mode: 'cors',
      method: 'POST', 
      headers: {
         Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
    .then(response => {
      if(response.status == 200){
          //Todo Switch screens and show user profile on new screen
      }
    })
    .catch((error) => {
      alert(error);
    });
  }


  render(){
    return (
      <View style={styles.container}>
        <Text style={styles.logo}>Fitness App</Text>
        <View style={styles.inputView}>
        <TextInput  
            style={styles.inputText}
            placeholder="username" 
            placeholderTextColor="#003f5c"
            onChangeText={text => this.setState({username:text})}/>
        </View>
        <View style={styles.inputView} >
          <TextInput  
            secureTextEntry
            style={styles.inputText}
            placeholder="Password" 
            placeholderTextColor="#003f5c"
            onChangeText={text => this.setState({password:text})}/>
        </View>
        <TouchableOpacity onPress={this.handleLoginOnPress} style={styles.loginBtn} >
          <Text style={styles.loginText}>LOGIN</Text>
        </TouchableOpacity>
        <TouchableOpacity onPress={this.props.navigation.navigate('SignupScreen')}  style={styles.loginBtn}>
          <Text style={styles.loginText}>Signup</Text>
        </TouchableOpacity>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignContent:'center',
    backgroundColor: '#003f5c',
    justifyContent: 'center',
  },
  logo:{
    fontWeight:"bold",
    fontSize:50,
    color:"#f55a42",
    marginBottom:40
  },
  inputView:{
    width:"80%",
    backgroundColor:"#465881",
    borderRadius:25,
    height:50,
    marginBottom:20,
    justifyContent:"center",
    padding:20
  },
  inputText:{
    height:50,
    color:"white"
  },
  loginText:{
    color:"white"
  },
  loginBtn:{
    width:"80%",
    backgroundColor:"#fb5b5a",
    borderRadius:25,
    height:50,
    alignItems:"center",
    justifyContent:"center",
    marginTop:40,
    marginBottom:10
  } 
});


