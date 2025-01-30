; Include necessary headers
!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "WinCore.nsh"
!include "Integration.nsh"

; Request application privileges for Windows Vista
RequestExecutionLevel user

; Build Unicode installer
Unicode True

; https://nsis.sourceforge.io/StrReplace_v4

!define Var0 $R0
!define Var1 $R1
!define Var2 $R2
!define Var3 $R3
!define Var4 $R4
!define Var5 $R5
!define Var6 $R6
!define Var7 $R7
!define Var8 $R8
 
!macro StrReplaceV4 Var Replace With In
 Push `${Replace}`
 Push `${With}`
 Push `${In}`
  Call StrReplaceV4
 Pop `${Var}`
!macroend
!define StrReplaceV4 `!insertmacro StrReplaceV4`
 
Function StrReplaceV4
Exch ${Var0} #in
Exch 1
Exch ${Var1} #with
Exch 2
Exch ${Var2} #replace
Push ${Var3}
Push ${Var4}
Push ${Var5}
Push ${Var6}
Push ${Var7}
Push ${Var8}
 
 StrCpy ${Var3} -1
 StrLen ${Var5} ${Var0}
 StrLen ${Var6} ${Var1}
 StrLen ${Var7} ${Var2}
 Loop:
  IntOp ${Var3} ${Var3} + 1
  StrCpy ${Var4} ${Var0} ${Var7} ${Var3}
  StrCmp ${Var3} ${Var5} End
  StrCmp ${Var4} ${Var2} 0 Loop
 
   StrCpy ${Var4} ${Var0} ${Var3}
   IntOp ${Var8} ${Var3} + ${Var7}
   StrCpy ${Var8} ${Var0} "" ${Var8}
   StrCpy ${Var0} ${Var4}${Var1}${Var8}
   IntOp ${Var3} ${Var3} + ${Var6}
   IntOp ${Var3} ${Var3} - 1
   IntOp ${Var5} ${Var5} - ${Var7}
   IntOp ${Var5} ${Var5} + ${Var6}
 
 Goto Loop
 End:
 
Pop ${Var8}
Pop ${Var7}
Pop ${Var6}
Pop ${Var5}
Pop ${Var4}
Pop ${Var3}
Pop ${Var2}
Exch
Pop ${Var1}
Exch ${Var0} #out
FunctionEnd
 
!undef Var8
!undef Var7
!undef Var6
!undef Var5
!undef Var4
!undef Var3
!undef Var2
!undef Var1
!undef Var0

; https://nsis.sourceforge.io/StrReplace_v4

; Variables
Var FIRST_NAME
Var LAST_NAME
Var TELEPHONE_NUMBER
Var TELEPHONE_URL
Var EMAIL_URL
Var OUTPUT_FILE_PATH

; Installer name and output
Name "Outlook Signature Installer"
OutFile "OutlookSignatureInstaller.exe"

; Define the source HTML file (to be included in the installer)
!define SOURCE_HTML "signature.htm"

; Define resource directory (to be included in the installer)
!define RESOURCE_DIR "signature-files"

; Define installation directory
InstallDir ""   ; Don't set a default $InstDir so we can detect /D= and InstallDirRegKey

; Pages
Page custom InputPage
Page custom TelephonePage
Page instfiles

Function .onInit
  SetShellVarContext Current

  ${If} $InstDir == "" ; No /D= nor InstallDirRegKey?
    GetKnownFolderPath $InstDir ${FOLDERID_RoamingAppData} ; This folder only exists on Win7+
    StrCmp $InstDir "" 0 +2 
  ${EndIf}

  
FunctionEnd

; Function to create a custom input page
Function InputPage
    nsDialogs::Create 1018
    ${If} $0 == error
        Abort
    ${EndIf}

    ; Label: Prompt for First Name
    ${NSD_CreateLabel} 0 0 100% 12u "Enter your First Name:"
    Pop $0

    ; Textbox for First Name
    ${NSD_CreateText} 0 12u 100% 12u ""
    Pop $1
    ${NSD_OnChange} $1 OnFirstNameChange

    ; Label: Prompt for Last Name
    ${NSD_CreateLabel} 0 30u 100% 12u "Enter your Last Name:"
    Pop $0

    ; Textbox for Last Name
    ${NSD_CreateText} 0 42u 100% 12u ""
    Pop $2
    ${NSD_OnChange} $2 OnLastNameChange

    ; Label: Prompt for E-Mail Address
    ${NSD_CreateLabel} 0 60u 100% 12u "Enter your email address:"
    Pop $0

    ; Textbox for E-Mail Address
    ${NSD_CreateText} 0 72u 100% 12u ""
    Pop $3
    ${NSD_OnChange} $3 OnEMailAddressChange
    
    ; Show the custom page
    nsDialogs::Show
FunctionEnd

Function OnFirstNameChange
    ${NSD_GetText} $1 $FIRST_NAME
FunctionEnd

Function OnLastNameChange
    ${NSD_GetText} $2 $LAST_NAME
FunctionEnd

Function OnEMailAddressChange
    ${NSD_GetText} $2 $EMAIL_URL
    ; Prefix with "mailto:"
    StrCpy $EMAIL_URL "mailto:$EMAIL_URL"
FunctionEnd

; Function to create the Telephone Page
Function TelephonePage
    nsDialogs::Create 1018
    ${If} $0 == error
        Abort
    ${EndIf}

    ; Label: Prompt for Telephone Number
    ${NSD_CreateLabel} 0 0 100% 12u "Enter your Telephone Number (e.g., +12 34567-890):"
    Pop $0

    ; Textbox for Telephone Number
    ${NSD_CreateText} 0 12u 100% 12u ""
    Pop $1
    ${NSD_OnChange} $1 OnTelephoneNumberChange

    ; Show the custom page
    nsDialogs::Show
FunctionEnd

Function OnTelephoneNumberChange
    ${NSD_GetText} $1 $TELEPHONE_NUMBER

    ; Initialize the URL variable
    StrCpy $TELEPHONE_URL ""

    ; Get the length of the telephone number string
    StrLen $3 $TELEPHONE_NUMBER

    ; Loop through each character of the telephone number
    ${For} $2 0 $3
        ; Extract the current character
        StrCpy $4 $TELEPHONE_NUMBER $2 1
        
        ; Check if the character is a valid digit (0-9)
        ${If} "$4" >= "0"
            ${AndIf} "$4" <= "9"
            StrCpy $TELEPHONE_URL "$TELEPHONE_URL$4" ; Append digit to URL
        ${ElseIf} "$4" == "+" ; Check if it's a '+'
            StrCpy $TELEPHONE_URL "$TELEPHONE_URL$4" ; Append '+' to URL
        ${EndIf}
    ${Next}

    ; Prefix the sanitized number with "tel:"
    StrCpy $TELEPHONE_URL "tel:$TELEPHONE_URL"
FunctionEnd

; Section to install and modify the HTML file
Section "Install Signature Files"

    ; Define output path
    StrCpy $OUTPUT_FILE_PATH "$APPDATA\Microsoft\Signatures\signature.htm"

    ; Create the installation directory
    CreateDirectory "$APPDATA\Microsoft\Signatures\signature-files"

    ; Extract the source HTML file (template) to a temporary location
    SetOutPath "$APPDATA\Microsoft\Signatures"
    File "${SOURCE_HTML}" ; Includes signature.htm in the installer
    StrCpy $0 "$InstDir\Microsoft\Signatures\signature-files\signature.htm"

    ; Open the signature file for reading
    FileOpen $1 $0 "r"
    FileOpen $2 $OUTPUT_FILE_PATH "w"

    ; Read and replace placeholders
    ${Do}
        FileRead $1 $3 ; Read a line into $3
        ${If} $3 == "" ; End of file
            ${Break}
        ${EndIf}

        ; Replace {{FIRST_NAME}}
		; ${StrReplaceV4} $Var "replace" "with" "in string"
		${StrReplaceV4} $3 "{{FIRST_NAME}}" $FIRST_NAME "$3"

        ; Replace {{LAST_NAME}}
        ${StrReplaceV4} $3 "{{LAST_NAME}}" $LAST_NAME "$3"

        ; Replace {{EMAIL_URL}}
		${StrReplaceV4} $3 "{{EMAIL_URL}}" $EMAIL_URL "$3"
        
        ; Replace {{TELEPHONE_URL}}
        ${StrReplaceV4} $3 "{{TELEPHONE_URL}}" $TELEPHONE_URL "$3"

        ; Replace {{TELEPHONE_NUMBER}}
        ${StrReplaceV4} $3 "{{TELEPHONE_NUMBER}}" $TELEPHONE_NUMBER "$3"

        ; Write the modified line to the final file
        FileWrite $2 $3
    ${Loop}

    ; Close files
    FileClose $1
    FileClose $2

    ; Delete the temporary template file
    Delete $0

    ; Install additional resource files
    SetOutPath "$APPDATA\Microsoft\Signatures\signature-files"
    ;SetOutPath "$InstDir\Microsoft\Signatures\signature-files"
    File /r "${RESOURCE_DIR}\*" ; Recursively includes all files and folders from the resources directory

    ; Confirmation message
    MessageBox MB_OK "The Outlook signature and its resources have been installed at: $APPDATA\Microsoft\Signatures"

; End of script
SectionEnd
