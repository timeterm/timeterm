\documentclass[11pt, letterpaper]{article}
\usepackage[a4paper, total={6in, 9in}]{geometry}
\usepackage{fontspec}
\setsansfont[BoldFont={Fira Sans}]{Fira Sans Light}
\setmonofont[Contextuals={Alternate}]{Fira Code}
\renewcommand{\familydefault}{\sfdefault}
\usepackage[dutch]{babel}
\usepackage{graphicx}
\usepackage{xcolor}
\usepackage{multicol}
\usepackage{titlesec}
\usepackage{minted}
\usepackage[xetex]{hyperref}
\usepackage{pdfpages}
\usepackage{etoolbox}
\usepackage{xpatch}
\usepackage{bookmark}

\PassOptionsToPackage{hyphens}{url}\usepackage{hyperref}

\renewcommand\theFancyVerbLine{\footnotesize\arabic{FancyVerbLine}}

\makeatletter
\AtBeginEnvironment{minted}{\dontdofcolorbox}
\def\dontdofcolorbox{\renewcommand\fcolorbox[4][]{##4}}
\xpatchcmd{\inputminted}{\minted@fvset}{\minted@fvset\dontdofcolorbox}{}{}
\xpatchcmd{\mintinline}{\minted@fvset}{\minted@fvset\dontdofcolorbox}{}{}
\makeatother

\title{PWS}
\author{Kees Blok, Rutger Broekhoff en Robert van der Maas}
\date{November 2020}

\begin{document}

\section{Broncode}

{{ range .Sources }}
	\subsection{
		\texorpdfstring{
			\nolinkurl{{ printf `{%s}` (texEscapeFull .Path) }}
		}{{ printf `{%s}` (texEscapeFull .Path) }}
	}
	
	\inputminted[
		fontsize=\footnotesize,
		breakanywhere,
		linenos,
		tabsize=2,
		breaklines,
		breaksymbolleft={\color{gray}\texttt{↳}},
		breaksymbolright={\color{gray}\texttt{⎦}},
	]{{ printf `{%s}` .Language }}{{ printf `{../../%s}` (texEscapeMintedPath .Path) }}
{{ end }}

\end{document}

