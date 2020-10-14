import React from 'react';
import { StyleSheet, Text,TouchableOpacity,TouchableHighlight,TextInput,View } from 'react-native';
import { string } from 'prop-types';
import { NavigationScreenProp } from 'react-navigation';
import { withNavigation } from 'react-navigation';
import { createStackNavigator } from '@react-navigation/stack';


export default class App extends React.Component {
    render(){
        return (
          <View style={styles.container}>
            <Text style={styles.logo}>Sign Up Here</Text>
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